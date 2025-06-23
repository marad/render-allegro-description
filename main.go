package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/a-h/templ"
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
	var filePath string
	var noPage bool
	flag.StringVar(&filePath, "file", "", "Ścieżka do pliku JSON (jeśli nie podano, wczytuje z stdin)")
	flag.StringVar(&filePath, "f", "", "Ścieżka do pliku JSON (skrót)")
	flag.BoolVar(&noPage, "no-page", false, "Renderuj tylko zawartość opisu bez pełnej strony HTML")
	flag.Parse()

	var jsonData string
	var err error

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
