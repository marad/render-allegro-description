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

func ParseItemList(data string) ([]Item, error) {
	var items []Item

	// Sprawdzenie czy dane nie są puste
	if len(data) == 0 {
		return items, fmt.Errorf("no JSON data provided")
	}

	// Parsowanie JSON z obsługą błędów
	err := json.Unmarshal([]byte(data), &items)
	if err != nil {
		return items, fmt.Errorf("invalid JSON: %v", err)
	}

	// Walidacja elementów
	for i, item := range items {
		if item.Type == "" {
			return items, fmt.Errorf("item %d has no type specified", i+1)
		}
		if item.Type == "IMAGE" && item.Url == "" {
			return items, fmt.Errorf("IMAGE item %d has no URL", i+1)
		}
		if item.Type == "TEXT" && item.Content == "" {
			return items, fmt.Errorf("TEXT item %d has no content", i+1)
		}
	}

	return items, nil
}

func ParseDescription(data string) (Description, error) {
	var desc Description

	// Sprawdzenie czy dane nie są puste
	if len(data) == 0 {
		return desc, fmt.Errorf("no JSON data provided")
	}

	// Parsowanie JSON z obsługą błędów
	err := json.Unmarshal([]byte(data), &desc)
	if err != nil {
		return desc, fmt.Errorf("invalid JSON: %v", err)
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
					return desc, fmt.Errorf("TEXT item in section %d has no content", i+1)
				}
			case "IMAGE":
				if item.Url == "" {
					return desc, fmt.Errorf("IMAGE item in section %d has no URL", i+1)
				}
			case "":
				return desc, fmt.Errorf("item in section %d has no type specified", i+1)
			default:
				// Nieznane typy są akceptowane, ale logujemy ostrzeżenie
				fmt.Printf("Warning: unknown item type '%s' in section %d\n", item.Type, i+1)
			}
		}
	}

	return desc, nil
}
