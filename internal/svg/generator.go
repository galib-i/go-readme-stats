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
)

//go:embed template.svg
var templateContent string

type SVGData struct {
	Theme         Theme
	Height        float64
	Header        string
	Languages     []stats.Lang // Includes colour codes
	LanguageCount int
}

// Generate creates an SVG of language statistics.
func Generate(theme, header string, languages []stats.Lang) (string, error) {
	languageCount := len(languages)

	data := SVGData{
		Theme:         GetTheme(theme),
		Height:        calculateSVGHeight(languageCount),
		Header:        header,
		Languages:     languages,
		LanguageCount: languageCount,
	}

	return generateSVG(data)
}

// calculateSVGHeight returns SVG height adjusted for language count.
func calculateSVGHeight(languageCount int) float64 {
	steps := (languageCount + 1) / 2
	return baseHeight + float64(steps-1)*heightStep
}

func generateSVG(data SVGData) (string, error) {
	tmpl, err := template.New(templateName).Funcs(template.FuncMap{
		"sumPrev": sumPreviousPercent,
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
