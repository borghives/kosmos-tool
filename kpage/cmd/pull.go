package cmd

import (
	"context"
	"fmt"
	"log"

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
		newIDStr, _ := cmd.Flags().GetString("new")
		mainIDStr, _ := cmd.Flags().GetString("main")

		if len(pageIDStrs) == 0 {
			log.Fatalf("Page ID is required")
		}

		var newID bson.ObjectID
		if newIDStr != "" {
			var err error
			newID, err = bson.ObjectIDFromHex(newIDStr)
			if err != nil {
				fmt.Printf("Failed to parse new page ID: newID=%s, error=%v\n", newIDStr, err)
			}
		}

		var mainID bson.ObjectID
		if mainIDStr != "" {
			var err error
			mainID, err = bson.ObjectIDFromHex(mainIDStr)
			if err != nil {
				log.Fatalf("Failed to parse main page ID: %v", err)
			}
		}

		pageIDs := make([]bson.ObjectID, 0, len(pageIDStrs)+2)
		for _, pageIdStr := range pageIDStrs {
			fmt.Printf("Action: Pulling page '%s'...\n", pageIdStr)

			id, err := bson.ObjectIDFromHex(pageIdStr)
			if err != nil {
				log.Fatalf("Failed to parse page ID: %v", err)
			}
			pageIDs = append(pageIDs, id)
		}

		if mainID != bson.NilObjectID {
			pageIDs = append(pageIDs, mainID)
		}

		if newID != bson.NilObjectID {
			pageIDs = append(pageIDs, newID)
		}

		pages, err := km.Filter[sitepages.SitePage](
			km.Fld("ID").ID().In(pageIDs...),
		).PullAll(context.Background())

		if err != nil {
			log.Fatalf("Failed to pull pages: %v", err)
		}

		for i, page := range pages {
			stanza, err := km.Filter[sitepages.Stanza](
				km.Fld("ID").ID().In(page.Contents...),
			).PullAll(context.Background())

			if err != nil {
				log.Fatalf("Failed to pull stanzas: %v", err)
			}

			pages[i].StanzaData = stanza

			if page.ID == newID {
				pages[i].Root = bson.NilObjectID
			}
		}

		//find main page and swap it to the front
		for i, page := range pages {
			if page.ID == mainID {
				remaining := append(pages[:i], pages[i+1:]...)
				pages = append([]sitepages.SitePage{page}, remaining...)
				break
			}
		}
		// Determine output
		if output == "" {
			return
		}

		err = sitepages.SaveSitePages(output, pages)
		if err != nil {
			log.Fatalf("Failed to write page to file: %v", err)
		}

		fmt.Printf("Successfully wrote page to '%s'\n", output)

	},
}

var pages []string

func init() {
	// Add the action to the context
	rootCmd.AddCommand(pullCmd)

	// Define flags
	pullCmd.Flags().StringSliceVarP(&pages, "page", "p", []string{}, "List of pages to pull (comma-separated or multiple flags)")
	pullCmd.Flags().StringP("new", "n", "", "ID of the page to be marked as a page genesis.")
	pullCmd.Flags().StringP("main", "m", "", "ID of the page to be marked as a main page.")
	pullCmd.Flags().StringP("output", "o", "", "Output file path (optional)")

	// Mark required flags
	pullCmd.MarkFlagRequired("page")
}
