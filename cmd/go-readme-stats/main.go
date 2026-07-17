package main

import (
	"context"
	_ "embed"
	"flag"
	"go-readme-stats/internal/stats"
	"go-readme-stats/internal/svg"
	"log"
	"os"
)

//go:embed ignored_languages.json
var ignoredLanguages []byte

func main() {
	theme := flag.String("theme", "dark", "theme for the card (dark, soft-dark or light)")
	header := flag.String("header", "Languages", "header text for the card")
	mode := flag.String("mode", "bytes", "percentage calculation mode (bytes or geometric)")
	flag.Parse()

	languages, err := stats.FetchStats(context.Background(), ignoredLanguages, *mode)
	if err != nil {
		log.Fatalf("Error fetching stats: %v", err)
	}

	svgContent, err := svg.Generate(*theme, *header, languages)
	if err != nil {
		log.Fatalf("Error generating SVG: %v", err)
	}

	err = os.WriteFile("languages.svg", []byte(svgContent), 0644)
	if err != nil {
		log.Fatalf("Error saving file: %v", err)
	}
}
