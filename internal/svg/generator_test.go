package svg

import (
	"go-readme-stats/internal/stats"
	"testing"
)

func TestCalculateSVGHeight(t *testing.T) {
	tests := []struct {
		name           string
		languageCount  int
		expectedHeight float64
	}{
		{"One language", 1, 114.5},
		{"Two languages", 2, 114.5},
		{"Three languages", 3, 134.5},
		{"Four languages", 4, 134.5},
		{"Five languages", 5, 154.5},
		{"Six languages", 6, 154.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateSVGHeight(tt.languageCount)
			if result != tt.expectedHeight {
				t.Errorf("calculateSVGHeight(%d) = %v, want %v", tt.languageCount, result, tt.expectedHeight)
			}
		})
	}
}

func TestSumPreviousPercent(t *testing.T) {
	languages := []stats.Lang{
		{Name: "Go", Percent: 45.5},
		{Name: "Java", Percent: 30.2},
		{Name: "JavaScript", Percent: 15.8},
		{Name: "Python", Percent: 8.5},
	}

	tests := []struct {
		name     string
		index    int
		expected float64
	}{
		{"First language", 0, 0.0},
		{"Second language", 1, 45.5},
		{"Third language", 2, 75.7},
		{"Fourth language", 3, 91.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sumPreviousPercent(languages, tt.index)
			if result != tt.expected {
				t.Errorf("sumPreviousPercent(languages, %d) = %v, want %v", tt.index, result, tt.expected)
			}
		})
	}
}
