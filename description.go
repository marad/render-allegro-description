package main

import (
	"encoding/json"
	"fmt"
)

type Item struct {
	Type    string `json:"type"`
	Url     string `json:"url"`
	Content string `json:"content"`
}

type Section struct {
	Items []Item `json:"items"`
}

type Description struct {
	Sections []Section `json:"sections"`
}

func ParseDescription(data string) (Description, error) {
	var desc Description

	// Sprawdzenie czy dane nie są puste
	if len(data) == 0 {
		return desc, fmt.Errorf("brak danych JSON")
	}

	// Parsowanie JSON z obsługą błędów
	err := json.Unmarshal([]byte(data), &desc)
	if err != nil {
		return desc, fmt.Errorf("nieprawidłowy JSON: %v", err)
	}

	// Walidacja podstawowej struktury
	if desc.Sections == nil {
		desc.Sections = []Section{}
	}

	// Walidacja elementów w sekcjach
	for i := range desc.Sections {
		if desc.Sections[i].Items == nil {
			desc.Sections[i].Items = []Item{}
		}

		// Walidacja typów elementów i ich zawartości
		for j := range desc.Sections[i].Items {
			item := &desc.Sections[i].Items[j]
			switch item.Type {
			case "TEXT":
				if item.Content == "" {
					return desc, fmt.Errorf("element TEXT w sekcji %d nie ma zawartości", i+1)
				}
			case "IMAGE":
				if item.Url == "" {
					return desc, fmt.Errorf("element IMAGE w sekcji %d nie ma URL", i+1)
				}
			case "":
				return desc, fmt.Errorf("element w sekcji %d nie ma określonego typu", i+1)
			default:
				// Nieznane typy są akceptowane, ale logujemy ostrzeżenie
				fmt.Printf("Ostrzeżenie: nieznany typ elementu '%s' w sekcji %d\n", item.Type, i+1)
			}
		}
	}

	return desc, nil
}
