package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
)

// RenderSd konwertuje JSON na HTML i zapisuje wynik do podanego writer'a
func RenderSd(jsonData string, out io.Writer) error {
	// Parsowanie JSON do struktury Description
	desc, err := ParseDescription(jsonData)
	if err != nil {
		return fmt.Errorf("błąd podczas parsowania JSON: %v", err)
	}

	// Renderowanie do HTML
	ctx := context.Background()
	component := descriptionPage(desc)

	err = component.Render(ctx, out)
	if err != nil {
		return fmt.Errorf("błąd podczas renderowania HTML: %v", err)
	}

	return nil
}

func main() {
	var filePath string
	flag.StringVar(&filePath, "file", "", "Ścieżka do pliku JSON (jeśli nie podano, wczytuje z stdin)")
	flag.StringVar(&filePath, "f", "", "Ścieżka do pliku JSON (skrót)")
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

	// Parsowanie JSON do struktury Description
	err = RenderSd(jsonData, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
