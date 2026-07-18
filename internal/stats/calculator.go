package stats

import (
	"cmp"
	"fmt"
	"math"
	"slices"
)

const (
	percentPrecision = 10 // One decimal precision (e.g. 10.4%)
)

// calculateStats computes language percentages by raw bytes or geometric mean.
// Languages are sorted by descending percentage, then ascending name.
// If more than maxLangs exist, languages at maxLangs-1 are grouped into "Other".
func calculateStats(languageTotals, languageFreq map[string]int, mode string, maxLangs int) []Lang {
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

	// Combine languages below top into "Other"
	if maxLangs > 0 && len(result) > maxLangs {
		top := maxLangs - 1
		var otherPercent float64

		for _, lang := range result[top:] {
			otherPercent += lang.Percent
		}

		result = append(result[:top], Lang{
			Name:    fmt.Sprintf("Other (%d)", len(result)-top),
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
