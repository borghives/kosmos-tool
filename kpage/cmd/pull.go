package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/borghives/kosmos-go"
	km "github.com/borghives/kosmos-go/model"
	"github.com/borghives/sitepages"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/spf13/cobra"
)

// Define the "pull" action command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull a page from a database",
	Run: func(cmd *cobra.Command, args []string) {
		pageIdStr, _ := cmd.Flags().GetString("page")
		output, _ := cmd.Flags().GetString("output")

		if pageIdStr == "" {
			log.Fatalf("Page ID is required")
		}

		fmt.Printf("Action: Pulling page '%s'...\n", pageIdStr)

		pageId, err := bson.ObjectIDFromHex(pageIdStr)
		if err != nil {
			log.Fatalf("Failed to parse page ID: %v", err)
		}
		page := kosmos.Filter[sitepages.SitePage](
			km.Fld("ID").Eq(pageId),
		).PullOne()

		if err != nil {
			log.Fatalf("Failed to pull page: %v", err)
		}

		// Determine output
		if output == "" {
			// Default to stdout
			fmt.Println("\n--- Page Content ---")
			fmt.Println(page.Title)
			fmt.Println("----------------------")
		} else {
			// Write to file
			err = os.WriteFile(output, []byte(page.Title), 0644)
			if err != nil {
				log.Fatalf("Failed to write page to file: %v", err)
			}
			fmt.Printf("Successfully wrote page to '%s'\n", output)
		}
	},
}

func init() {
	// Add the action to the context
	rootCmd.AddCommand(pullCmd)

	// Define flags
	pullCmd.Flags().StringP("page", "p", "", "Page name to pull")
	pullCmd.Flags().StringP("output", "o", "", "Output file path (optional)")

	// Mark required flags
	pullCmd.MarkFlagRequired("page")
}
