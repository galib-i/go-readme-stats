package stats

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCallAPI_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"test": "data"}`))
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
}

func TestCallAPI_InvalidURL(t *testing.T) {
	_, err := callAPI(context.Background(), "invalid-url")
	if err == nil {
		t.Error("callAPI() expected error for invalid URL")
	}
}
