package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/a-h/templ"
	"github.com/docopt/docopt-go"
)

func RenderSd(jsonData string, out io.Writer, fullPage bool) error {
	desc, err := ParseDescription(jsonData)
	if err != nil {
		return fmt.Errorf("błąd podczas parsowania JSON: %v", err)
	}

	ctx := context.Background()
	var component templ.Component

	if fullPage {
		component = descriptionPage(desc)
	} else {
		component = description(desc)
	}

	err = component.Render(ctx, out)
	if err != nil {
		return fmt.Errorf("błąd podczas renderowania HTML: %v", err)
	}

	return nil
}

func main() {
	usage := `Render SD - narzędzie do renderowania opisów na podstawie JSON

Usage:
  render-sd [--file=<path> | -f <path>] [--no-page]
  render-sd -h | --help

Options:
  -f <path>, --file=<path>  Ścieżka do pliku JSON (jeśli nie podano, wczytuje z stdin)
  --no-page                 Renderuj tylko zawartość opisu bez pełnej strony HTML
  -h, --help                Pokaż tę pomoc`

	opts, err := docopt.ParseDoc(usage)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Błąd podczas parsowania argumentów: %v\n", err)
		os.Exit(1)
	}

	filePath, _ := opts.String("--file")
	noPage, _ := opts.Bool("--no-page")

	var jsonData string

	if filePath != "" {
		// Wczytywanie z pliku
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Błąd podczas wczytywania pliku: %v\n", err)
			os.Exit(1)
		}
		jsonData = string(fileContent)
	} else {
		// Wczytywanie ze stdin
		stdinContent, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Błąd podczas wczytywania ze stdin: %v\n", err)
			os.Exit(1)
		}
		jsonData = string(stdinContent)
	}

	// Renderowanie JSON do HTML
	err = RenderSd(jsonData, os.Stdout, !noPage)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
