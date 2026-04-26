package cmd

import (
	"os"

	"github.com/borghives/kosmos-go"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ksecret",
	Short: "A CLI tool to manage secret",
}

// Execute is called by main.main().
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(extractCmd)
	rootCmd.AddCommand(rotateCmd)

	rootCmd.PersistentFlags().StringP("project", "p", "", "Project ID")

	cobra.OnInitialize(func() {
		kosmos.IgniteBase(rootCmd, "tool.env")
	})
}
