package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/borghives/kosmos-go"
	"github.com/borghives/kosmos-go/observation"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"gopkg.in/yaml.v3"
)

// IndexKey represents a single field and its direction in a MongoDB index.
type IndexKey struct {
	Field string `yaml:"field"`
	Order int32  `yaml:"order"`
}

// IndexConfig represents the configuration for a single MongoDB index.
type IndexConfig struct {
	Keys   []IndexKey `yaml:"keys"`
	Unique bool       `yaml:"unique"`
}

type TimeseriesOption struct {
	TimeField   string `yaml:"time-field"`
	MetaField   string `yaml:"meta-field"`
	Granularity string `yaml:"granularity"`
}

func (t TimeseriesOption) ToMongoOptions() *options.TimeSeriesOptionsBuilder {
	fmt.Printf("    -> Timeseries options: time-field: %s, meta-field: %s, granularity: %s\n", t.TimeField, t.MetaField, t.Granularity)
	return options.TimeSeries().
		SetTimeField(t.TimeField).
		SetMetaField(t.MetaField).
		SetGranularity(t.Granularity)
}

// CollectionConfig represents the configuration for a MongoDB collection and its internal indexes.
type CollectionConfig struct {
	Name       string            `yaml:"name"`
	Indexes    []IndexConfig     `yaml:"indexes"`
	Timeseries *TimeseriesOption `yaml:"timeseries,omitempty"`
}

// DatabaseConfig represents a MongoDB database and its collections.
type DatabaseConfig struct {
	Name        string             `yaml:"name"`
	Collections []CollectionConfig `yaml:"collections"`
}

type UserConfig struct {
	Name      string   `yaml:"name"`
	SecretSrc string   `yaml:"secret-src"`
	Read      []string `yaml:"read"`
	ReadWrite []string `yaml:"write"`
	Creator   []string `yaml:"creator"`
	Admin     bool     `yaml:"admin"`
}

// DataConfig is the root structure for data initialization YAML file.
type DataConfig struct {
	Databases []DatabaseConfig `yaml:"databases"`
	Users     []UserConfig     `yaml:"users"`
}

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Manage MongoDB data, collections, and indexes",
}

var dbDeclareCmd = &cobra.Command{
	Use:   "declare",
	Short: "Declare database, collections, and indexes from a YAML configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		filePath, _ := cmd.Flags().GetString("file")
		if filePath == "" {
			log.Fatalf("Error: --file flag is required")
		}

		fmt.Printf("Action: Declaring data from '%s'...\n", filePath)

		// 1. Read YAML configuration
		yamlFile, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("Failed to read file: %v", err)
		}

		var dataConfig DataConfig
		err = yaml.Unmarshal(yamlFile, &dataConfig)
		if err != nil {
			log.Fatalf("Failed to parse YAML: %v", err)
		}

		// 3. Process the configuration
		for _, userConfig := range dataConfig.Users {
			processingUser(userConfig)
		}

		for _, dbConfig := range dataConfig.Databases {
			processingDatabase(dbConfig)
		}

		fmt.Println("Data initialization completed successfully.")
	},
}

func createCollection(ctx context.Context, db *mongo.Database, collConfig CollectionConfig) {
	fmt.Printf("  Processing collection: %s\n", collConfig.Name)

	// Optional: explicitly create collection if it doesn't exist
	fmt.Printf("    -> Calling db.CreateCollection for %s...\n", collConfig.Name)
	collOptions := options.CreateCollection()
	if collConfig.Timeseries != nil {
		fmt.Printf("    -> timeseries options for %s...\n", collConfig.Name)
		collOptions = collOptions.SetTimeSeriesOptions(collConfig.Timeseries.ToMongoOptions())
	}

	err := db.CreateCollection(ctx, collConfig.Name, collOptions)
	if err != nil {
		fmt.Printf("    <- db.CreateCollection returned: %v\n", err)
	} else {
		fmt.Printf("    <- db.CreateCollection succeeded\n")
	}

	collection := db.Collection(collConfig.Name)

	// Create indexes
	for _, indexConfig := range collConfig.Indexes {
		if len(indexConfig.Keys) == 0 {
			continue
		}

		// Build index keys preserving order
		keysDoc := bson.D{}
		for _, k := range indexConfig.Keys {
			keysDoc = append(keysDoc, bson.E{Key: k.Field, Value: k.Order})
		}

		indexModel := mongo.IndexModel{
			Keys: keysDoc,
		}

		if indexConfig.Unique {
			indexModel.Options = options.Index().SetUnique(true)
		}

		fmt.Printf("    -> Calling CreateOne for index %v...\n", keysDoc)
		name, err := collection.Indexes().CreateOne(ctx, indexModel)
		if err != nil {
			fmt.Printf("    <- CreateOne returned name: %s, error: %v\n", name, err)
		} else {
			fmt.Printf("    <- CreateOne name: %s succeeded\n", name)
		}
		if err != nil {
			log.Fatalf("    Failed to create index on collection %s: %v", collConfig.Name, err)
		}
		fmt.Printf("    Created index: %s\n", name)
	}
}

func processingDatabase(dbConfig DatabaseConfig) {

	branchName := dbConfig.Name
	if branchName == "" {
		branchName = "main"
	}

	fmt.Printf("Processing database branch name: %s\n", branchName)
	dataverse := kosmos.SummonObservationFor(observation.PurposeAffinityCreator)
	defer dataverse.Close()

	// ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	// defer cancel()

	ctx := context.Background()

	db := dataverse.BranchDatabase(dbConfig.Name)

	for _, collConfig := range dbConfig.Collections {
		createCollection(ctx, db, collConfig)
	}
}

func processingUser(userConfig UserConfig) {
	fmt.Printf("Processing user: %s\n", userConfig.Name)
	dataverse := kosmos.SummonObservationFor(observation.PurposeAffinityAdmin)
	defer dataverse.Close()

	password, err := kosmos.CollapseSecret(userConfig.SecretSrc)
	if err != nil {
		log.Fatalf("Failed to extract password source : %v", err)
	}

	responsibility := observation.MemberResponsibility{
		ReadDbs:      userConfig.Read,
		ReadWriteDbs: userConfig.ReadWrite,
		CreatorDbs:   userConfig.Creator,
		IsAdmin:      userConfig.Admin,
	}

	err = dataverse.UpdateMember(userConfig.Name, password, responsibility, true)
	if err != nil {
		log.Fatalf("Failed to set user: %v", err)
	}
}

func init() {
	dbCmd.AddCommand(dbDeclareCmd)

	// Define persistent flags or local flags
	dbDeclareCmd.Flags().StringP("file", "f", "", "Path to the data initialization YAML file")
	dbDeclareCmd.MarkFlagRequired("file")
}
