package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/goccy/go-yaml"
)

const (
	sourceURL   = "https://raw.githubusercontent.com/github-linguist/linguist/refs/heads/main/lib/linguist/languages.yml"
	outputPath  = "internal/stats/colours.json"
	httpTimeout = 30 * time.Second

	// A partial parse would wipe colours from the card, so refuse to write one.
	minLanguages = 100
)

func main() {
	if err := fetchLanguageColours(); err != nil {
		log.Fatalf("Failed to fetch language colours: %v", err)
	}

	fmt.Println("Successfully fetched language colours.")
}

// fetchLanguageColours downloads GitHub's language definitions and keeps only
// the colours.
func fetchLanguageColours() error {
	client := &http.Client{Timeout: httpTimeout}

	resp, err := client.Get(sourceURL)
	if err != nil {
		return fmt.Errorf("failed to fetch YAML: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch YAML: unexpected status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var definitions map[string]any
	if err := yaml.Unmarshal(body, &definitions); err != nil {
		return fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	colours := make(map[string]string)

	for language, definition := range definitions {
		fields, isDefinition := definition.(map[string]any)
		if !isDefinition {
			continue
		}

		// "color" is GitHub's own key name.
		if colour, hasColour := fields["color"].(string); hasColour {
			colours[language] = colour
		}
	}

	if len(colours) < minLanguages {
		return fmt.Errorf("refusing to write %s: only %d languages parsed", outputPath, len(colours))
	}

	jsonData, err := json.MarshalIndent(colours, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return os.WriteFile(outputPath, jsonData, 0o644)
}
