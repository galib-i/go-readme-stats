package stats

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func useTestServer(t *testing.T, url string) {
	t.Helper()

	previous := apiBaseURL
	apiBaseURL = url
	t.Cleanup(func() { apiBaseURL = previous })
}

func TestCallAPI_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"test": "data"}`))
	}))
	defer server.Close()

	body, err := callAPI(context.Background(), server.URL)
	if err != nil {
		t.Errorf("callAPI() error = %v", err)
	}

	expected := `{"test": "data"}`
	if string(body) != expected {
		t.Errorf("callAPI() = %s, want %s", string(body), expected)
	}
}

func TestCallAPI_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	_, err := callAPI(context.Background(), server.URL)
	if err == nil {
		t.Error("callAPI() expected error for 404 response")
	}

	if !isStatus(err, http.StatusNotFound) {
		t.Errorf("callAPI() error = %v, want a 404 status error", err)
	}
}

func TestCallAPI_RetriesServerError(t *testing.T) {
	calls := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls == 1 {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"test": "data"}`))
	}))
	defer server.Close()

	body, err := callAPI(context.Background(), server.URL)
	if err != nil {
		t.Fatalf("callAPI() error = %v", err)
	}

	if calls != 2 {
		t.Errorf("callAPI() made %d requests, want 2", calls)
	}

	if string(body) != `{"test": "data"}` {
		t.Errorf("callAPI() = %s, want %s", string(body), `{"test": "data"}`)
	}
}

func TestCallAPI_InvalidURL(t *testing.T) {
	_, err := callAPI(context.Background(), "invalid-url")
	if err == nil {
		t.Error("callAPI() expected error for invalid URL")
	}
}

func TestFetchRepositories_PaginatesOwnedRepos(t *testing.T) {
	requestedPages := []string{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		if affiliation := query.Get("affiliation"); affiliation != "owner" {
			t.Errorf("affiliation = %q, want %q", affiliation, "owner")
		}

		page := query.Get("page")
		requestedPages = append(requestedPages, page)

		// The first page is full, the second is short, which ends the loop.
		batch := make([]repository, 0, reposPerPage)
		if page == "1" {
			for i := range reposPerPage {
				batch = append(batch, repository{Name: fmt.Sprintf("repo-%d", i)})
			}
		} else {
			batch = append(batch, repository{Name: "last-repo"})
		}

		if err := json.NewEncoder(w).Encode(batch); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	useTestServer(t, server.URL)

	repos, err := fetchRepositories(context.Background())
	if err != nil {
		t.Fatalf("fetchRepositories() error = %v", err)
	}

	if len(repos) != reposPerPage+1 {
		t.Errorf("fetchRepositories() returned %d repos, want %d", len(repos), reposPerPage+1)
	}

	if len(requestedPages) != 2 {
		t.Fatalf("fetchRepositories() requested pages %v, want 2 pages", requestedPages)
	}
}

// newGitHubStub serves the three endpoints FetchStats uses.
func newGitHubStub(t *testing.T, repos []repository, languages map[string]map[string]int) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/user":
			_, _ = w.Write([]byte(`{"login": "octocat"}`))
		case "/user/repos":
			if err := json.NewEncoder(w).Encode(repos); err != nil {
				t.Errorf("failed to encode repos: %v", err)
			}
		default:
			repoName, found := strings.CutPrefix(r.URL.Path, "/repos/octocat/")
			if !found {
				http.NotFound(w, r)

				return
			}

			langs, known := languages[strings.TrimSuffix(repoName, "/languages")]
			if !known {
				w.WriteHeader(http.StatusNotFound)

				return
			}

			if err := json.NewEncoder(w).Encode(langs); err != nil {
				t.Errorf("failed to encode languages: %v", err)
			}
		}
	}))
}

func TestFetchStats_SkipsForksAndIgnoredLanguages(t *testing.T) {
	server := newGitHubStub(t,
		[]repository{{Name: "active"}, {Name: "forked", Fork: true}, {Name: "empty"}},
		map[string]map[string]int{
			"active": {"Go": 1000, "HTML": 500},
		},
	)
	defer server.Close()

	useTestServer(t, server.URL)

	languages, err := FetchStats(context.Background(), []string{"HTML"}, "bytes", 6)
	if err != nil {
		t.Fatalf("FetchStats() error = %v", err)
	}

	if len(languages) != 1 {
		t.Fatalf("FetchStats() returned %d languages, want 1: %+v", len(languages), languages)
	}

	if languages[0].Name != "Go" || languages[0].Percent != 100 {
		t.Errorf("FetchStats() = %+v, want Go at 100%%", languages[0])
	}

	if languages[0].Colour != "#00ADD8" {
		t.Errorf("Colour = %s, want #00ADD8", languages[0].Colour)
	}
}

func TestFetchStats_FailsOnRepoError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/user":
			_, _ = w.Write([]byte(`{"login": "octocat"}`))
		case "/user/repos":
			_, _ = w.Write([]byte(`[{"name": "broken"}]`))
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	useTestServer(t, server.URL)

	if _, err := FetchStats(context.Background(), nil, "bytes", 6); err == nil {
		t.Error("FetchStats() expected an error when a repository lookup fails")
	}
}

func TestFetchStats_FailsWhenNothingIsFound(t *testing.T) {
	// A 404 means "no language data", so every repo here is skipped.
	server := newGitHubStub(t, []repository{{Name: "empty"}}, nil)
	defer server.Close()

	useTestServer(t, server.URL)

	if _, err := FetchStats(context.Background(), nil, "bytes", 6); err == nil {
		t.Error("FetchStats() expected an error when no languages are found")
	}
}
