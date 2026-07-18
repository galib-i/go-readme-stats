package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/goccy/go-yaml"
)

const (
	url        = "https://raw.githubusercontent.com/github-linguist/linguist/refs/heads/main/lib/linguist/languages.yml"
	outputPath = "internal/stats/colours.json"
)

func main() {
	if err := fetchLanguageColours(); err != nil {
		log.Fatalf("Failed to fetch language colours: %v", err)
	}
	fmt.Println("Successfully fetched language colours.")
}

// fetchLanguageColours downloads and converts GitHub's language colours to JSON.
// Fetches the official YAML file and extracts only the colour mappings.
func fetchLanguageColours() error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch YAML: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	var data map[string]any
	if err := yaml.Unmarshal(body, &data); err != nil {
		return fmt.Errorf("failed to unmarshal YAML: %v", err)
	}

	filtered := make(map[string]string)
	for k, v := range data {
		if m, ok := v.(map[string]any); ok {
			if colour, ok := m["color"].(string); ok {
				filtered[k] = colour
			}
		}
	}

	jsonData, err := json.MarshalIndent(filtered, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	return os.WriteFile(outputPath, jsonData, 0644)
}
