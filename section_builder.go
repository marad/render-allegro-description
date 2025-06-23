package main

func BuildSections(items []Item) []Section {
	if len(items) == 0 {
		return []Section{}
	}

	sections := []Section{}
	currentSection := Section{Items: []Item{}}

	for i := 0; i < len(items); i++ {
		currentItem := items[i]

		// Dodaj aktualny item do bieżącej sekcji
		currentSection.Items = append(currentSection.Items, currentItem)

		// Sprawdź czy należy zakończyć bieżącą sekcję
		shouldEndSection := false

		if i == len(items)-1 || len(currentSection.Items) >= 2 {
			// To jest ostatni item, lub sekcja zawiera już 2 elementy - kończymy sekcję
			shouldEndSection = true
		} else {
			nextItem := items[i+1]
			// Jeśli aktualny item to TEXT i następny item to też TEXT, kończymy sekcję
			if currentItem.Type == "TEXT" && nextItem.Type == "TEXT" {
				shouldEndSection = true
			}
		}

		if shouldEndSection {
			sections = append(sections, currentSection)
			currentSection = Section{Items: []Item{}}
		}
	}

	return sections
}
