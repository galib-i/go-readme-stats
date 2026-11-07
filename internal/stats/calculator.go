package stats

import (
	"cmp"
	"fmt"
	"math"
	"slices"
)

const (
	percentScale = 10 // Rounds to one decimal place, e.g. 10.4%
)

// scoreLanguages ranks languages by their share. It leaves percentages
// unrounded for FormatStats to sum, and ignores max_langs so one fetch can
// serve several cards.
func scoreLanguages(languageTotals, languageFreq map[string]int, mode string) []Lang {
	scores := make(map[string]float64)
	var totalScore float64

	for lang, bytes := range languageTotals {
		freq := languageFreq[lang]
		var score float64

		switch mode {
		case "geometric": // Geometric mean: sqrt(bytes * freq)
			score = math.Sqrt(float64(bytes) * float64(freq))
		default:
			score = float64(bytes)
		}

		scores[lang] = score
		totalScore += score
	}

	result := make([]Lang, 0, len(scores))
	for lang, score := range scores {
		result = append(result, Lang{
			Name:    lang,
			Percent: score / totalScore * 100,
		})
	}

	// Highest percentage first; ties go by name, so the order is stable.
	slices.SortFunc(result, func(a, b Lang) int {
		if c := cmp.Compare(b.Percent, a.Percent); c != 0 {
			return c
		}
		return cmp.Compare(a.Name, b.Name)
	})

	return result
}

// FormatStats groups everything past the top maxLangs-1 into "Other", then
// rounds the percentages. It does not modify the input, so cards can share it.
func FormatStats(scored []Lang, maxLangs int) []Lang {
	result := slices.Clone(scored)

	if maxLangs > 0 && len(result) > maxLangs {
		top := maxLangs - 1
		var otherPercent float64

		for _, lang := range result[top:] {
			otherPercent += lang.Percent
		}

		result = append(result[:top], Lang{
			Name:    fmt.Sprintf("Other (%d)", len(result)-top),
			Percent: otherPercent,
			Colour:  defaultColour,
		})
	}

	for i := range result {
		result[i].Percent = roundPercent(result[i].Percent)
	}

	// Rounded rows can add up to 99.9 or 100.1, so the last row takes the
	// remainder and the bar always ends at exactly 100%.
	if len(result) > 0 {
		var rounded float64
		for _, lang := range result[:len(result)-1] {
			rounded += lang.Percent
		}
		result[len(result)-1].Percent = roundPercent(100 - rounded)
	}

	return result
}

func roundPercent(percent float64) float64 {
	return math.Round(percent*percentScale) / percentScale
}
