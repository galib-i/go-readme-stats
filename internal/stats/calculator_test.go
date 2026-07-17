package stats

import (
	"testing"
)

func TestCalculateStats_Raw(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]int
		expected []Lang
	}{
		{
			name: "Basic calculation",
			input: map[string]int{
				"Go":         5000,
				"JavaScript": 3000,
				"Python":     2000,
			},
			expected: []Lang{
				{Name: "Go", Percent: 50.0},
				{Name: "JavaScript", Percent: 30.0},
				{Name: "Python", Percent: 20.0},
			},
		},
		{
			name: "Single language",
			input: map[string]int{
				"Go": 1000,
			},
			expected: []Lang{
				{Name: "Go", Percent: 100.0},
			},
		},
		{
			name:     "Empty input",
			input:    map[string]int{},
			expected: []Lang{},
		},
		{
			name: "Many languages (should group into Other)",
			input: map[string]int{
				"Go":         5000,
				"JavaScript": 2000,
				"Python":     1500,
				"Java":       1000,
				"C++":        800,
				"Rust":       400,
				"TypeScript": 200,
				"HTML":       100,
			},
			expected: []Lang{
				{Name: "Go", Percent: 45.5},
				{Name: "JavaScript", Percent: 18.2},
				{Name: "Python", Percent: 13.6},
				{Name: "Java", Percent: 9.1},
				{Name: "C++", Percent: 7.3},
				{Name: "Other (3)", Percent: 6.4},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create frequency map with 1 for all languages
			freq := make(map[string]int, len(tt.input))
			for lang := range tt.input {
				freq[lang] = 1
			}

			result := calculateStats(tt.input, freq, "")

			if len(result) != len(tt.expected) {
				t.Errorf("got %d languages, expected %d", len(result), len(tt.expected))
				return
			}

			for i, lang := range result {
				if lang.Name != tt.expected[i].Name {
					t.Errorf("[%d] Name = %s, expected %s", i, lang.Name, tt.expected[i].Name)
				}
				if lang.Percent != tt.expected[i].Percent {
					t.Errorf("[%d] Percent = %v, expected %v", i, lang.Percent, tt.expected[i].Percent)
				}
			}
		})
	}
}

func TestCalculateStats_Geometric(t *testing.T) {
	tests := []struct {
		name     string
		totals   map[string]int
		freq     map[string]int
		expected []Lang
	}{
		{
			name: "Basic geometric",
			totals: map[string]int{
				"Go":         5000,
				"JavaScript": 3000,
			},
			freq: map[string]int{
				"Go":         2,
				"JavaScript": 8,
			},
			expected: []Lang{
				{Name: "JavaScript", Percent: 60.8},
				{Name: "Go", Percent: 39.2},
			},
		},
		{
			name: "Same percentage uses name sort",
			totals: map[string]int{
				"B": 1000,
				"A": 2000,
				"C": 1000,
			},
			freq: map[string]int{
				"A": 1,
				"B": 1,
				"C": 1,
			},
			expected: []Lang{
				{Name: "A", Percent: 41.4},
				{Name: "B", Percent: 29.3},
				{Name: "C", Percent: 29.3},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateStats(tt.totals, tt.freq, "geometric")

			if len(result) != len(tt.expected) {
				t.Errorf("got %d languages, expected %d", len(result), len(tt.expected))
				return
			}

			for i, lang := range result {
				if lang.Name != tt.expected[i].Name {
					t.Errorf("[%d] Name = %s, expected %s", i, lang.Name, tt.expected[i].Name)
				}
				if lang.Percent != tt.expected[i].Percent {
					t.Errorf("[%d] Percent = %v, expected %v", i, lang.Percent, tt.expected[i].Percent)
				}
			}
		})
	}
}

func TestRoundPercent(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected float64
	}{
		{"Round down", 10.14, 10.1},
		{"Round up", 10.16, 10.2},
		{"Exact", 10.1, 10.1},
		{"Zero", 0.0, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := roundPercent(tt.input)
			if result != tt.expected {
				t.Errorf("roundPercent(%v) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}
