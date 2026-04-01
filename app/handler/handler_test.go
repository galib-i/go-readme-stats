package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-readme-stats/app/stats"

	"github.com/gin-gonic/gin"
)

// Mock implementations for tests
func init() {
	// Default success mock
	FetchStats = func(ctx context.Context, ignoredLanguagesData []byte, mode string) ([]stats.Lang, error) {
		return []stats.Lang{
			{Name: "Go", Percent: 45.5},
			{Name: "Java", Percent: 30.2},
			{Name: "JavaScript", Percent: 15.8},
			{Name: "Python", Percent: 8.5},
		}, nil
	}

	GenerateSVG = func(theme, header string, languages []stats.Lang) (string, error) {
		return "<svg>mock</svg>", nil
	}
}

func TestGetLanguageStats_Success(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "Default parameters",
			url:      "/langs",
			expected: "<svg>mock</svg>",
		},
		{
			name:     "Custom theme",
			url:      "/langs?theme=light",
			expected: "<svg>mock</svg>",
		},
		{
			name:     "Custom header",
			url:      "/langs?header=My%20Languages",
			expected: "<svg>mock</svg>",
		},
		{
			name:     "Both parameters",
			url:      "/langs?theme=dark&header=Code%20Stats",
			expected: "<svg>mock</svg>",
		},
	}

	gin.SetMode(gin.TestMode)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.GET("/langs", GetLanguageStats)

			req, _ := http.NewRequest("GET", tt.url, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
			}

			if w.Header().Get("Content-Type") != "image/svg+xml" {
				t.Errorf("Expected Content-Type 'image/svg+xml', got '%s'", w.Header().Get("Content-Type"))
			}

			if body := w.Body.String(); body != tt.expected {
				t.Errorf("Expected body '%s', got '%s'", tt.expected, body)
			}
		})
	}
}

func TestGetLanguageStats_StatsFetchFailure(t *testing.T) {
	originalFetch := FetchStats
	FetchStats = func(ctx context.Context, ignoredLanguagesData []byte, mode string) ([]stats.Lang, error) {
		return nil, errors.New("API rate limit exceeded")
	}
	defer func() { FetchStats = originalFetch }()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/langs", GetLanguageStats)

	req, _ := http.NewRequest("GET", "/langs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	if body := w.Body.String(); body != "Error fetching stats" {
		t.Errorf("Expected error message, got '%s'", body)
	}
}

func TestGetLanguageStats_SVGFailure(t *testing.T) {
	originalGenerate := GenerateSVG
	GenerateSVG = func(theme, header string, languages []stats.Lang) (string, error) {
		return "", errors.New("template error")
	}
	defer func() { GenerateSVG = originalGenerate }()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/langs", GetLanguageStats)

	req, _ := http.NewRequest("GET", "/langs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	if body := w.Body.String(); body != "Error generating SVG" {
		t.Errorf("Expected error message, got '%s'", body)
	}
}

func TestGetLanguageStats_InvalidTheme(t *testing.T) {
	originalGenerate := GenerateSVG
	GenerateSVG = func(theme, header string, languages []stats.Lang) (string, error) {
		if theme == "invalid" {
			return "", errors.New("invalid theme")
		}
		return "<svg>mock</svg>", nil
	}
	defer func() { GenerateSVG = originalGenerate }()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/langs", GetLanguageStats)

	req, _ := http.NewRequest("GET", "/langs?theme=invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestGetLanguageStats_EmptyLanguages(t *testing.T) {
	originalFetch := FetchStats
	FetchStats = func(ctx context.Context, ignoredLanguagesData []byte, mode string) ([]stats.Lang, error) {
		return []stats.Lang{}, nil
	}
	defer func() { FetchStats = originalFetch }()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/langs", GetLanguageStats)

	req, _ := http.NewRequest("GET", "/langs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	if body := w.Body.String(); body != "<svg>mock</svg>" {
		t.Errorf("Expected SVG output even with empty languages")
	}
}
