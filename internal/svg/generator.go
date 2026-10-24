package svg

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"

	"go-readme-stats/internal/stats"
)

const (
	baseHeight   = 114.5
	heightStep   = 20.0
	templateName = "template.svg"
	maxNameLen   = 12
)

//go:embed template.svg
var templateContent string

type SVGData struct {
	Theme     Theme
	Height    float64
	Header    string
	Languages []stats.Lang // Includes colour codes
}

// Generate creates an SVG of language statistics.
func Generate(theme, header string, languages []stats.Lang) (string, error) {
	data := SVGData{
		Theme:     GetTheme(theme),
		Height:    calculateSVGHeight(len(languages)),
		Header:    header,
		Languages: languages,
	}

	return generateSVG(data)
}

// calculateSVGHeight grows the card by one row per pair of languages, since the
// legend shows two per row.
func calculateSVGHeight(languageCount int) float64 {
	rows := max(1, (languageCount+1)/2)

	return baseHeight + float64(rows-1)*heightStep
}

// truncateName shortens long names to fit the legend. It counts letters rather
// than bytes, so accented names are never cut in half.
func truncateName(name string) string {
	runes := []rune(name)
	if len(runes) <= maxNameLen {
		return name
	}

	return string(runes[:maxNameLen]) + "..."
}

func generateSVG(data SVGData) (string, error) {
	tmpl, err := template.New(templateName).Funcs(template.FuncMap{
		"sumPrev":  sumPreviousPercent,
		"add":      func(a, b int) int { return a + b },
		"mul":      func(a, b int) int { return a * b },
		"div":      func(a, b int) int { return a / b },
		"mod":      func(a, b int) int { return a % b },
		"truncate": truncateName,
	}).Parse(templateContent)

	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// sumPreviousPercent calculates cumulative percentage for stacked progress bars.
func sumPreviousPercent(languages []stats.Lang, idx int) float64 {
	sum := 0.0
	for i := range idx {
		sum += languages[i].Percent
	}
	return sum
}
