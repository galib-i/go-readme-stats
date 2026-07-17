package stats

import (
	"cmp"
	"fmt"
	"math"
	"slices"
)

const (
	maxVisibleLanguages = 6
	topLanguagesCount   = maxVisibleLanguages - 1 // Number of individual languages before grouping into "Other"
	percentPrecision    = 10                      // One decimal precision (e.g. 10.4%)
)

// calculateStats computes language percentages by raw bytes or geometric mean.
// Languages are sorted by descending percentage, then ascending name.
// If more than maxVisibleLanguages exist, languages beyond topLanguagesCount are grouped into "Other".
func calculateStats(languageTotals, languageFreq map[string]int, mode string) []Lang {
	scores := make(map[string]float64)
	var totalScore float64

	for lang, bytes := range languageTotals {
		freq := languageFreq[lang]
		var score float64

		switch mode {
		case "geometric": // Geometric mean: sqrt(bytes * freq)
			score = math.Sqrt(float64(bytes) * float64(freq))
		default: // Raw byte count
			score = float64(bytes)
		}

		scores[lang] = score
		totalScore += score
	}

	var result []Lang
	for lang, score := range scores {
		result = append(result, Lang{
			Name:    lang,
			Percent: score / totalScore * 100,
		})
	}

	// Sort by percentage (desc), then by name (asc) for consistent ordering
	slices.SortFunc(result, func(a, b Lang) int {
		if c := cmp.Compare(b.Percent, a.Percent); c != 0 {
			return c
		}
		return cmp.Compare(a.Name, b.Name)
	})

	// Combine languages below topLanguagesCount into "Other"
	if len(result) > maxVisibleLanguages {
		var otherPercent float64
		for _, lang := range result[topLanguagesCount:] {
			otherPercent += lang.Percent
		}

		result = append(result[:topLanguagesCount], Lang{
			Name:    fmt.Sprintf("Other (%d)", len(result)-topLanguagesCount),
			Percent: otherPercent,
		})
	}

	for i := range result {
		result[i].Percent = roundPercent(result[i].Percent)
	}

	return result
}

func roundPercent(percent float64) float64 {
	return math.Round(percent*percentPrecision) / percentPrecision
}
