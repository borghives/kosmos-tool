package cmd

import (
	"os"

	"github.com/borghives/kosmos-go"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kpage",
	Short: "A CLI tool to manage page",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(pullCmd)

	cobra.OnInitialize(func() {
		rootCmd.PersistentFlags().StringP("uri", "u", "", "MongoDB connection URI")
		kosmos.Ignite(rootCmd, "tool.env")
	})
}
