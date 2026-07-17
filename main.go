package main

import (
	"context"
	_ "embed"
	"go-readme-stats/app/stats"
	"go-readme-stats/app/svg"
	"log"
	"os"
)

//go:embed ignored_languages.json
var ignoredLanguages []byte

func main() {
	theme := svg.DefaultTheme // Dark
	header := "Languages"
	mode := "bytes"

	languages, err := stats.FetchStats(context.Background(), ignoredLanguages, mode)
	if err != nil {
		log.Fatalf("Error fetching stats: %v", err)
	}

	svgContent, err := svg.Generate(theme, header, languages)
	if err != nil {
		log.Fatalf("Error generating SVG: %v", err)
	}

	err = os.WriteFile("languages.svg", []byte(svgContent), 0644)
	if err != nil {
		log.Fatalf("Error saving file: %v", err)
	}
}
