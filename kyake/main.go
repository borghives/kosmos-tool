package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"

	"github.com/joho/godotenv"
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

func main() {
	envFile := flag.String("env", "", "Path to the .env file")
	tplFile := flag.String("in", "", "Path to the template file")
	outFile := flag.String("out", "", "Path to the output file (optional, defaults to stdout)")
	
	flag.Parse()

	if *envFile == "" || *tplFile == "" {
		fmt.Fprintln(os.Stderr, "Usage: envsubst -env <file> -in <template> [-out <output>]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	envMap, err := godotenv.Read(*envFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading env file: %v\n", err)
		os.Exit(1)
	}

	tplBytes, err := os.ReadFile(*tplFile)
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
	if *outFile == "" {
		out = os.Stdout
	} else {
		// Ensure directory exists
		if err := os.MkdirAll(filepath.Dir(*outFile), 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
			os.Exit(1)
		}
		f, err := os.Create(*outFile)
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
}
