package stats

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"slices"
	"sync"
)

//go:embed colours.json
var coloursJSON []byte

const (
	defaultColour = "#F0F6FC"
)

type repository struct {
	Name string `json:"name"`
	Fork bool   `json:"fork"`
}

type Lang struct {
	Name    string
	Percent float64
	Colour  string
}

// FetchStats retrieves language statistics for the authenticated user.
// Excludes forked repositories and languages in the ignored list.
func FetchStats(ctx context.Context, ignoredLanguages []string, mode string, maxLangs int) ([]Lang, error) {
	repos, err := fetchRepoNames(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repositories: %w", err)
	}

	username, err := getUsername(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get authenticated user: %w", err)
	}

	languageTotals := make(map[string]int)
	languageFreq := make(map[string]int)

	var wg sync.WaitGroup
	var mu sync.Mutex

	maxConcurrentRequests := 20
	sem := make(chan struct{}, maxConcurrentRequests)

	for _, repo := range repos {
		if repo.Fork {
			continue
		}

		wg.Go(func() {

			sem <- struct{}{}
			defer func() { <-sem }()

			languages, err := fetchRepoLanguages(ctx, username, repo.Name)
			if err != nil {
				log.Printf("Warning: Failed to fetch languages for %s: %v", repo.Name, err)
				return
			}

			mu.Lock()
			defer mu.Unlock()

			for lang, bytes := range languages {
				if slices.Contains(ignoredLanguages, lang) {
					continue
				}

				languageTotals[lang] += bytes
				languageFreq[lang]++
			}
		})
	}

	wg.Wait()

	stats := calculateStats(languageTotals, languageFreq, mode, maxLangs)
	if err := addLanguageColours(stats); err != nil {
		return nil, fmt.Errorf("failed to add colours: %w", err)
	}
	return stats, nil
}

func addLanguageColours(languages []Lang) error {
	colours, err := loadLanguageColours()
	if err != nil {
		log.Printf("Warning: Failed to load colours: %v", err)
		colours = make(map[string]string)
	}

	for i := range languages {
		colour, exists := colours[languages[i].Name]
		if !exists {
			colour = defaultColour
		}
		languages[i].Colour = colour
	}

	return nil
}

func loadLanguageColours() (map[string]string, error) {
	var colours map[string]string
	if err := json.Unmarshal(coloursJSON, &colours); err != nil {
		return nil, fmt.Errorf("failed to parse embedded colours: %w", err)
	}

	return colours, nil
}
