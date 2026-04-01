package cmd

import (
	"os"

	"github.com/borghives/kosmos-go/ether"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kpage",
	Short: "A CLI tool to manage pages",
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
		ether.CollapseConstants().MergeFromFile("tool.env").MergeFromCmd(rootCmd)
		ether.CollapseObserverConstants().MergeFromFile("tool.env").MergeFromCmd(rootCmd)
	})
}
