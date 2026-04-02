package cmd

import (
	"fmt"
	"log"
	"os"

	km "github.com/borghives/kosmos-go"
	"github.com/borghives/sitepages"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/spf13/cobra"
)

// Define the "pull" action command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull a page from a database",
	Run: func(cmd *cobra.Command, args []string) {
		pageIDStrs, _ := cmd.Flags().GetStringSlice("page")
		output, _ := cmd.Flags().GetString("output")

		if len(pageIDStrs) == 0 {
			log.Fatalf("Page ID is required")
		}

		pageIDs := make([]any, 0, len(pageIDStrs))
		for _, pageIdStr := range pageIDStrs {
			fmt.Printf("Action: Pulling page '%s'...\n", pageIdStr)

			id, err := bson.ObjectIDFromHex(pageIdStr)
			if err != nil {
				log.Fatalf("Failed to parse page ID: %v", err)
			}
			pageIDs = append(pageIDs, id)
		}

		pages, err := km.Filter[sitepages.SitePage](
			km.Fld("ID").In(pageIDs...),
		).PullAll()

		if err != nil {
			log.Fatalf("Failed to pull pages: %v", err)
		}

		for _, page := range pages {

			// Determine output
			if output == "" {
				// Default to stdout
				fmt.Println("\n--- Page Title ---")
				fmt.Println(page.Title)
				fmt.Println("----------------------")
			} else {
				// Write to file
				err := os.WriteFile(output, []byte(page.Title), 0644)
				if err != nil {
					log.Fatalf("Failed to write page to file: %v", err)
				}

				fmt.Printf("Successfully wrote page to '%s'\n", output)
			}
		}
	},
}

var pages []string

func init() {
	// Add the action to the context
	rootCmd.AddCommand(pullCmd)

	// Define flags
	pullCmd.Flags().StringSliceVarP(&pages, "page", "p", []string{}, "List of pages to pull (comma-separated or multiple flags)")
	pullCmd.Flags().StringP("output", "o", "", "Output file path (optional)")

	// Mark required flags
	pullCmd.MarkFlagRequired("page")
}
