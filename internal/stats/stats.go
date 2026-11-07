package stats

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"slices"
	"sync"
)

//go:embed colours.json
var coloursJSON []byte

const (
	defaultColour = "#F0F6FC"

	maxConcurrentRequests = 20
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

// FetchStats returns language statistics for the user's own repositories,
// excluding forks and ignored languages.
//
// Every language comes back unrounded and untruncated, so call FormatStats for
// one card's rows. It fails outright rather than return partial numbers, so a
// bad token or rate limit can never publish a misleading card.
func FetchStats(ctx context.Context, ignoredLanguages []string, mode string) ([]Lang, error) {
	repos, err := fetchRepositories(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repositories: %w", err)
	}

	username, err := fetchUsername(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get authenticated user: %w", err)
	}

	var (
		languageTotals = make(map[string]int)
		languageFreq   = make(map[string]int)
		failures       []error
		mu             sync.Mutex
		wg             sync.WaitGroup
	)

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
				// A repository with no commits has no languages, so a 404 is not a failure.
				if isStatus(err, http.StatusNotFound) {
					log.Printf("Warning: Skipping %s: repository has no language data", repo.Name)

					return
				}

				mu.Lock()
				failures = append(failures, fmt.Errorf("%s: %w", repo.Name, err))
				mu.Unlock()

				return
			}

			mu.Lock()
			defer mu.Unlock()

			for lang, size := range languages {
				if slices.Contains(ignoredLanguages, lang) {
					continue
				}

				languageTotals[lang] += size
				languageFreq[lang]++
			}
		})
	}

	wg.Wait()

	if len(failures) > 0 {
		return nil, fmt.Errorf("%d of %d repositories failed (first error: %w)", len(failures), len(repos), failures[0])
	}

	if len(languageTotals) == 0 {
		return nil, errors.New("no languages found: check the token scopes and the ignore list")
	}

	languages := scoreLanguages(languageTotals, languageFreq, mode)
	addLanguageColours(languages)

	return languages, nil
}

// addLanguageColours sets each colour, using the default for unknown languages.
func addLanguageColours(languages []Lang) {
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
}

func loadLanguageColours() (map[string]string, error) {
	var colours map[string]string
	if err := json.Unmarshal(coloursJSON, &colours); err != nil {
		return nil, fmt.Errorf("failed to parse embedded colours: %w", err)
	}

	return colours, nil
}
