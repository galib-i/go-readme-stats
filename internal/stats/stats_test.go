package stats

import (
	"testing"
)

func TestLoadLanguageColours(t *testing.T) {
	colours, err := loadLanguageColours()
	if err != nil {
		t.Errorf("loadLanguageColours() error = %v", err)
		return
	}

	if len(colours) == 0 {
		t.Error("loadLanguageColours() returned empty map")
	}

	if goColour, exists := colours["Go"]; !exists {
		t.Error("Go language colour not found")
	} else if goColour != "#00ADD8" {
		t.Errorf("Go colour = %s, want #00ADD8", goColour)
	}
}

func TestParseIgnoredLanguages(t *testing.T) {
	testData := []byte(`["HTML", "CSS"]`)
	ignored, err := parseIgnoredLanguages(testData)
	if err != nil {
		t.Errorf("parseIgnoredLanguages() error = %v", err)
		return
	}

	if _, exists := ignored["HTML"]; !exists {
		t.Error("HTML should be in ignored languages")
	}

	if _, exists := ignored["CSS"]; !exists {
		t.Error("CSS should be in ignored languages")
	}
}

func TestParseIgnoredLanguages_InvalidFile(t *testing.T) {
	_, err := parseIgnoredLanguages([]byte("{invalid}"))
	if err == nil {
		t.Error("parseIgnoredLanguages() expected error for invalid JSON")
	}
}

func TestAddLanguageColours(t *testing.T) {
	languages := []Lang{
		{Name: "Go", Percent: 50.0},
		{Name: "UnknownLanguage", Percent: 30.0},
		{Name: "Java", Percent: 20.0},
	}

	err := addLanguageColours(languages)
	if err != nil {
		t.Errorf("addLanguageColours() error = %v", err)
		return
	}

	if languages[0].Colour != "#00ADD8" {
		t.Errorf("Go language colour = %s, want #00ADD8", languages[0].Colour)
	}

	if languages[1].Colour != defaultColour {
		t.Errorf("Unknown language colour = %s, want %s", languages[1].Colour, defaultColour)
	}

	if languages[2].Colour == "" {
		t.Error("Java language should have a colour assigned")
	}
}
