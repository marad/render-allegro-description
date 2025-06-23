package main

import (
	"strings"
	"testing"
)

func TestParseItemList_ValidData(t *testing.T) {
	jsonData := `[{"type":"TEXT","content":"Hello"},{"type":"IMAGE","url":"http://example.com/image.jpg"}]`
	items, err := ParseItemList(jsonData)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if len(items) != 2 {
		t.Errorf("Expected 2 items, got: %d", len(items))
	}
	if items[0].Type != "TEXT" || items[0].Content != "Hello" {
		t.Errorf("First item not parsed correctly: %+v", items[0])
	}
	if items[1].Type != "IMAGE" || items[1].Url != "http://example.com/image.jpg" {
		t.Errorf("Second item not parsed correctly: %+v", items[1])
	}
}

func TestParseItemList_EmptyData(t *testing.T) {
	_, err := ParseItemList("")
	if err == nil || err.Error() != "no JSON data provided" {
		t.Errorf("Expected 'no JSON data provided' error, got: %v", err)
	}
}

func TestParseItemList_InvalidJSON(t *testing.T) {
	_, err := ParseItemList("not a json")
	if err == nil || !strings.HasPrefix(err.Error(), "invalid JSON") {
		t.Errorf("Expected 'invalid JSON' error, got: %v", err)
	}
}

func TestParseItemList_MissingType(t *testing.T) {
	jsonData := `[{"content":"Hello"}]`
	_, err := ParseItemList(jsonData)
	if err == nil || err.Error() != "item 1 has no type specified" {
		t.Errorf("Expected 'item 1 has no type specified' error, got: %v", err)
	}
}

func TestParseItemList_ImageNoURL(t *testing.T) {
	jsonData := `[{"type":"IMAGE"}]`
	_, err := ParseItemList(jsonData)
	if err == nil || err.Error() != "IMAGE item 1 has no URL" {
		t.Errorf("Expected 'IMAGE item 1 has no URL' error, got: %v", err)
	}
}

func TestParseItemList_TextNoContent(t *testing.T) {
	jsonData := `[{"type":"TEXT"}]`
	_, err := ParseItemList(jsonData)
	if err == nil || err.Error() != "TEXT item 1 has no content" {
		t.Errorf("Expected 'TEXT item 1 has no content' error, got: %v", err)
	}
}
