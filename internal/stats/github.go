package stats

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func fetchRepoNames(ctx context.Context) ([]repository, error) {
	url := "https://api.github.com/user/repos"
	body, err := callAPI(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repos: %w", err)
	}

	var repos []repository
	if err = json.Unmarshal(body, &repos); err != nil {
		return nil, fmt.Errorf("failed to unmarshal repos: %w", err)
	}

	return repos, nil
}

func fetchRepoLanguages(ctx context.Context, username, repoName string) (map[string]int, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/languages", username, repoName)
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

func getUsername(ctx context.Context) (string, error) {
	body, err := callAPI(ctx, "https://api.github.com/user")
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

// callAPI makes authenticated HTTP requests to the GitHub API.
// Uses GITHUB_TOKEN environment variable for authentication if available.
func callAPI(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "go-readme-stats")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request to %s: %w", url, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status %d for URL %s", resp.StatusCode, url)
	}

	return io.ReadAll(resp.Body)

}
