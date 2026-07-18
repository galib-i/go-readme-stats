package main

import (
	"context"
	"log"
	"os"

	"go-readme-stats/internal/stats"
	"go-readme-stats/internal/svg"

	"github.com/goccy/go-yaml"
)

type cardConfig struct {
	Output   string   `yaml:"output"`
	Theme    string   `yaml:"theme"`
	Header   string   `yaml:"header"`
	Mode     string   `yaml:"mode"`
	MaxLangs int      `yaml:"max_langs"`
	Ignore   []string `yaml:"ignore"`
}

type config struct {
	Cards []cardConfig `yaml:"cards"`
}

// loadConfig reads stats.yml and returns the parsed configuration.
func loadConfig(path string) config {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read %s: %v", path, err)
	}

	var cfg config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("Failed to parse %s: %v", path, err)
	}

	// Apply per-card defaults for any omitted fields
	for i := range cfg.Cards {
		c := &cfg.Cards[i]
		if c.Output == "" {
			c.Output = "languages.svg"
		}
		if c.Theme == "" {
			c.Theme = "dark"
		}
		if c.Header == "" {
			c.Header = "Languages"
		}
		if c.Mode == "" {
			c.Mode = "bytes"
		}
		if c.MaxLangs == 0 {
			c.MaxLangs = 6
		}
	}

	return cfg
}

func main() {
	cfg := loadConfig("stats.yml")

	for _, card := range cfg.Cards {
		languages, err := stats.FetchStats(context.Background(), card.Ignore, card.Mode, card.MaxLangs)
		if err != nil {
			log.Fatalf("Error fetching stats for %s: %v", card.Output, err)
		}

		svgContent, err := svg.Generate(card.Theme, card.Header, languages)
		if err != nil {
			log.Fatalf("Error generating SVG for %s: %v", card.Output, err)
		}

		if err = os.WriteFile(card.Output, []byte(svgContent), 0644); err != nil {
			log.Fatalf("Error saving %s: %v", card.Output, err)
		}

		log.Printf("Generated %s", card.Output)
	}
}
