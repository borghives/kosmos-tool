package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"text/template"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// ProcessTemplate substitutes environment variables from envMap into the template.
func ProcessTemplate(tpl string, envMap map[string]string) (string, error) {
	t, err := template.New("config").Option("missingkey=zero").Parse(tpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, envMap)
	return buf.String(), err
}

var rootCmd = &cobra.Command{
	Use:   "kyake",
	Short: "A CLI tool to make YAML files from templates and env files",
	Run: func(cmd *cobra.Command, args []string) {
		envFiles, err := cmd.Flags().GetStringSlice("env")
		if err != nil {
			log.Fatalf("Failed to get environment files: %v", err)
		}
		tplFile, err := cmd.Flags().GetString("in")
		if err != nil {
			log.Fatalf("Failed to get input template file: %v", err)
		}
		outFile, err := cmd.Flags().GetString("out")
		if err != nil {
			log.Fatalf("Failed to get output file: %v", err)
		}

		if tplFile == "" {
			log.Fatalf("No input YAML template file provided.")
		}

		envMap, err := godotenv.Read(envFiles...)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading env file: %v\n", err)
			os.Exit(1)
		}

		tplBytes, err := os.ReadFile(tplFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading template file: %v\n", err)
			os.Exit(1)
		}

		result, err := ProcessTemplate(string(tplBytes), envMap)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error processing template: %v\n", err)
			os.Exit(1)
		}

		var out io.Writer
		if outFile == "" {
			out = os.Stdout
		} else {
			f, err := os.Create(outFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
				os.Exit(1)
			}
			defer f.Close()
			out = f
		}

		if _, err := fmt.Fprint(out, result); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var env []string

func init() {
	rootCmd.PersistentFlags().StringSliceVarP(&env, "env", "e", []string{}, "Environment file")
	rootCmd.PersistentFlags().StringP("in", "i", "", "Input YAML template file")
	rootCmd.PersistentFlags().StringP("out", "o", "", "Output YAML file")
}
