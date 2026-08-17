package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoutes(t *testing.T) {
	s := New()
	ts := httptest.NewServer(s.Routes())
	defer ts.Close()

	tests := []struct {
		name       string
		path       string
		wantStatus int
	}{
		{"healthz", "/healthz", http.StatusOK},
		{"readyz when ready", "/readyz", http.StatusOK},
		{"version", "/version", http.StatusOK},
		{"ping", "/api/v1/ping", http.StatusOK},
		{"unknown route", "/does-not-exist", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.Get(ts.URL + tt.path)
			if err != nil {
				t.Fatalf("GET %s: %v", tt.path, err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("GET %s: got status %d, want %d", tt.path, resp.StatusCode, tt.wantStatus)
			}
		})
	}
}

func TestReadyzReflectsDraining(t *testing.T) {
	s := New()
	ts := httptest.NewServer(s.Routes())
	defer ts.Close()

	s.SetNotReady()

	resp, err := http.Get(ts.URL + "/readyz")
	if err != nil {
		t.Fatalf("GET /readyz: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("got status %d, want %d after SetNotReady", resp.StatusCode, http.StatusServiceUnavailable)
	}

	var body map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body["status"] != "draining" {
		t.Errorf("got status field %q, want %q", body["status"], "draining")
	}
}

func TestVersionPayload(t *testing.T) {
	s := New()
	ts := httptest.NewServer(s.Routes())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/version")
	if err != nil {
		t.Fatalf("GET /version: %v", err)
	}
	defer resp.Body.Close()

	var body struct {
		Version   string `json:"version"`
		GoVersion string `json:"goVersion"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.GoVersion == "" {
		t.Error("expected goVersion to be populated")
	}
}
