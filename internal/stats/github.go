package stats

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// apiBaseURL can be overridden so tests can point at a local server.
var apiBaseURL = "https://api.github.com"

const (
	reposPerPage   = 100 // GitHub's maximum, so a short page means the end
	requestTimeout = 30 * time.Second
	maxAttempts    = 2 // one retry for temporary failures
)

// httpClient is shared so connections are reused. The timeout stops a stalled
// request from hanging the workflow.
var httpClient = &http.Client{Timeout: requestTimeout}

type apiError struct {
	StatusCode int
	URL        string
}

func (e *apiError) Error() string {
	msg := fmt.Sprintf("HTTP request failed with status %d for URL %s", e.StatusCode, e.URL)

	switch e.StatusCode {
	case http.StatusUnauthorized:
		msg += ": the token is missing or invalid"
	case http.StatusForbidden, http.StatusTooManyRequests:
		msg += ": rate limited, or the token lacks the required scopes"
	}

	return msg
}

func isStatus(err error, code int) bool {
	var apiErr *apiError
	return errors.As(err, &apiErr) && apiErr.StatusCode == code
}

// fetchRepositories lists the user's own repositories. Forks, and repositories
// they only collaborate on, are left out.
func fetchRepositories(ctx context.Context) ([]repository, error) {
	var repos []repository

	for page := 1; ; page++ {
		url := fmt.Sprintf("%s/user/repos?affiliation=owner&per_page=%d&page=%d", apiBaseURL, reposPerPage, page)

		body, err := callAPI(ctx, url)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch repos (page %d): %w", page, err)
		}

		var batch []repository
		if err := json.Unmarshal(body, &batch); err != nil {
			return nil, fmt.Errorf("failed to unmarshal repos: %w", err)
		}

		repos = append(repos, batch...)

		if len(batch) < reposPerPage {
			return repos, nil
		}
	}
}

func fetchRepoLanguages(ctx context.Context, username, repoName string) (map[string]int, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/languages", apiBaseURL, username, repoName)
	body, err := callAPI(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch languages for %s: %w", repoName, err)
	}

	var languages map[string]int
	if err := json.Unmarshal(body, &languages); err != nil {
		return nil, fmt.Errorf("failed to unmarshal languages for %s: %w", repoName, err)
	}

	return languages, nil
}

func fetchUsername(ctx context.Context) (string, error) {
	body, err := callAPI(ctx, apiBaseURL+"/user")
	if err != nil {
		return "", fmt.Errorf("failed to get user info: %w", err)
	}

	var user struct {
		Login string `json:"login"`
	}

	if err := json.Unmarshal(body, &user); err != nil {
		return "", fmt.Errorf("failed to parse user info: %w", err)
	}

	return user.Login, nil
}

// callAPI calls the GitHub API with the token from GITHUB_TOKEN, retrying once
// when a failure looks temporary.
func callAPI(ctx context.Context, url string) ([]byte, error) {
	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		body, retryable, err := doRequest(ctx, url)
		if err == nil {
			return body, nil
		}

		lastErr = err
		if !retryable || attempt == maxAttempts {
			break
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(time.Duration(attempt) * time.Second):
		}
	}

	return nil, lastErr
}

// doRequest sends one request. The bool reports whether the failure is worth
// retrying.
func doRequest(ctx context.Context, url string) ([]byte, bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, false, fmt.Errorf("failed to create request: %w", err)
	}

	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "go-readme-stats")

	resp, err := httpClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, false, ctx.Err()
		}

		return nil, true, fmt.Errorf("failed to make HTTP request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Read the body out so the connection can be reused.
		_, _ = io.Copy(io.Discard, resp.Body)

		// Only server errors are worth another try; a bad token or a rate limit
		// needs a person to fix it.
		return nil, resp.StatusCode >= 500, &apiError{StatusCode: resp.StatusCode, URL: url}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, false, fmt.Errorf("failed to read response from %s: %w", url, err)
	}

	return body, false, nil
}
