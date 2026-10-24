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

// calculateStats computes language percentages by raw bytes or geometric mean.
// Languages are sorted by descending percentage, then ascending name.
// Anything past the top maxLangs-1 is grouped into "Other".
func calculateStats(languageTotals, languageFreq map[string]int, mode string, maxLangs int) []Lang {
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

	var result []Lang
	for lang, score := range scores {
		result = append(result, Lang{
			Name:    lang,
			Percent: score / totalScore * 100,
		})
	}

	// Highest percentage first; the name breaks ties so the order never wobbles.
	slices.SortFunc(result, func(a, b Lang) int {
		if c := cmp.Compare(b.Percent, a.Percent); c != 0 {
			return c
		}
		return cmp.Compare(a.Name, b.Name)
	})

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
	return math.Round(percent*percentScale) / percentScale
}
