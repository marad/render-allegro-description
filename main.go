package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/a-h/templ"
	"github.com/docopt/docopt-go"
)

const usage = `Render SD - tool for rendering descriptions based on JSON

Usage:
  render-sd [--file=<path> | -f <path>] [--no-page]
  render-sd -h | --help

Options:
  -f <path>, --file=<path>  Path to JSON file (if not provided, reads from stdin)
  --no-page                 Render only description content without full HTML page
  -h, --help                Show this help`

func main() {
	opts, err := docopt.ParseDoc(usage)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing arguments: %v\n", err)
		os.Exit(1)
	}

	filePath, _ := opts.String("--file")
	noPage, _ := opts.Bool("--no-page")

	var jsonData string

	if filePath != "" {
		// Reading from file
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}
		jsonData = string(fileContent)
	} else {
		// Reading from stdin
		stdinContent, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
			os.Exit(1)
		}
		jsonData = string(stdinContent)
	}

	// Reading description from JSON
	desc, err := ReadDescription(jsonData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading description: %v\n", err)
		os.Exit(1)
	}

	// Rendering JSON to HTML
	err = RenderSd(desc, os.Stdout, !noPage)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func ReadDescription(jsonData string) (Description, error) {
	desc, err := ParseDescription(jsonData)
	if err != nil {
		return Description{}, fmt.Errorf("error parsing JSON: %v", err)
	}
	return desc, nil
}

func RenderSd(desc Description, out io.Writer, fullPage bool) error {
	ctx := context.Background()
	var component templ.Component

	if fullPage {
		component = descriptionPage(desc)
	} else {
		component = description(desc)
	}

	err := component.Render(ctx, out)
	if err != nil {
		return fmt.Errorf("error rendering HTML: %v", err)
	}

	return nil
}
