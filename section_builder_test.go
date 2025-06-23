package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildSections_EmptyInput(t *testing.T) {
	items := []Item{}

	sections := BuildSections(items)

	assert.Empty(t, sections)
}

func TestBuildSections_SingleTextItem(t *testing.T) {
	items := []Item{
		{Type: "TEXT", Content: "Single text"},
	}

	sections := BuildSections(items)

	assert.Len(t, sections, 1)
	assert.Len(t, sections[0].Items, 1)
	assert.Equal(t, "TEXT", sections[0].Items[0].Type)
	assert.Equal(t, "Single text", sections[0].Items[0].Content)
}

func TestBuildSections_TwoConsecutiveTextItems(t *testing.T) {
	items := []Item{
		{Type: "TEXT", Content: "First text"},
		{Type: "TEXT", Content: "Second text"},
	}

	sections := BuildSections(items)

	// Kolejne dwa TEXT itemy powinny być w osobnych sekcjach
	assert.Len(t, sections, 2)
	assert.Len(t, sections[0].Items, 1)
	assert.Equal(t, "TEXT", sections[0].Items[0].Type)
	assert.Equal(t, "First text", sections[0].Items[0].Content)

	assert.Len(t, sections[1].Items, 1)
	assert.Equal(t, "TEXT", sections[1].Items[0].Type)
	assert.Equal(t, "Second text", sections[1].Items[0].Content)
}

func TestBuildSections_TextAndImageTogether(t *testing.T) {
	items := []Item{
		{Type: "TEXT", Content: "Description"},
		{Type: "IMAGE", Url: "image.jpg"},
	}

	sections := BuildSections(items)

	// TEXT i IMAGE powinny być w tej samej sekcji
	assert.Len(t, sections, 1)
	assert.Len(t, sections[0].Items, 2)
	assert.Equal(t, "TEXT", sections[0].Items[0].Type)
	assert.Equal(t, "Description", sections[0].Items[0].Content)
	assert.Equal(t, "IMAGE", sections[0].Items[1].Type)
	assert.Equal(t, "image.jpg", sections[0].Items[1].Url)
}

func TestBuildSections_ExampleFromRequirements(t *testing.T) {
	// [text-a, image-f, text-b, text-c, image-d, text-e]
	// Oczekiwane sekcje:
	// - text-a, image-f
	// - text-b
	// - text-c, image-d
	// - text-e
	items := []Item{
		{Type: "TEXT", Content: "text-a"},
		{Type: "IMAGE", Content: "image-f"},
		{Type: "TEXT", Content: "text-b"},
		{Type: "TEXT", Content: "text-c"},
		{Type: "IMAGE", Content: "image-d"},
		{Type: "TEXT", Content: "text-e"},
	}

	sections := BuildSections(items)

	assert.Len(t, sections, 4)

	// Sekcja 1: text-a, image-f
	assert.Len(t, sections[0].Items, 2)
	assert.Equal(t, "TEXT", sections[0].Items[0].Type)
	assert.Equal(t, "text-a", sections[0].Items[0].Content)
	assert.Equal(t, "IMAGE", sections[0].Items[1].Type)
	assert.Equal(t, "image-f", sections[0].Items[1].Content)

	// Sekcja 2: text-b (zakończona bo następny też TEXT)
	assert.Len(t, sections[1].Items, 1)
	assert.Equal(t, "TEXT", sections[1].Items[0].Type)
	assert.Equal(t, "text-b", sections[1].Items[0].Content)

	// Sekcja 3: text-c, image-d (maksymalnie 2 elementy)
	assert.Len(t, sections[2].Items, 2)
	assert.Equal(t, "TEXT", sections[2].Items[0].Type)
	assert.Equal(t, "text-c", sections[2].Items[0].Content)
	assert.Equal(t, "IMAGE", sections[2].Items[1].Type)
	assert.Equal(t, "image-d", sections[2].Items[1].Content)

	// Sekcja 4: text-e
	assert.Len(t, sections[3].Items, 1)
	assert.Equal(t, "TEXT", sections[3].Items[0].Type)
	assert.Equal(t, "text-e", sections[3].Items[0].Content)
}

func TestBuildSections_ThreeConsecutiveTextItems(t *testing.T) {
	items := []Item{
		{Type: "TEXT", Content: "First"},
		{Type: "TEXT", Content: "Second"},
		{Type: "TEXT", Content: "Third"},
	}

	sections := BuildSections(items)

	// Każdy TEXT po TEXT powinien być w osobnej sekcji
	assert.Len(t, sections, 3)
	assert.Len(t, sections[0].Items, 1)
	assert.Equal(t, "First", sections[0].Items[0].Content)
	assert.Len(t, sections[1].Items, 1)
	assert.Equal(t, "Second", sections[1].Items[0].Content)
	assert.Len(t, sections[2].Items, 1)
	assert.Equal(t, "Third", sections[2].Items[0].Content)
}

func TestBuildSections_MixedTypes(t *testing.T) {
	items := []Item{
		{Type: "IMAGE", Content: "image1"},
		{Type: "VIDEO", Content: "video1"},
		{Type: "TEXT", Content: "text1"},
		{Type: "LINK", Content: "link1"},
	}

	sections := BuildSections(items)

	// Żadne z nich nie są kolejnymi TEXT, więc wszystko w jednej sekcji
	// ale sekcja ma już maksymalnie 2 elementy, więc będą 2 sekcje
	assert.Len(t, sections, 2)
	assert.Len(t, sections[0].Items, 2)
	assert.Equal(t, "IMAGE", sections[0].Items[0].Type)
	assert.Equal(t, "VIDEO", sections[0].Items[1].Type)
	assert.Len(t, sections[1].Items, 2)
	assert.Equal(t, "TEXT", sections[1].Items[0].Type)
	assert.Equal(t, "LINK", sections[1].Items[1].Type)
}
