package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	srv := New(8080, "", "output")
	req := httptest.NewRequest("GET", "/api/health", nil)
	w := httptest.NewRecorder()

	srv.handleHealth(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp HealthResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Status != "ok" {
		t.Fatalf("expected status ok, got %s", resp.Status)
	}
	if resp.Version != version {
		t.Fatalf("expected version %s, got %s", version, resp.Version)
	}
}

func TestAuthMiddleware_NoKey(t *testing.T) {
	srv := New(8080, "", "output")
	called := false
	handler := srv.authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	if !called {
		t.Fatal("handler should be called when no API key is configured")
	}
}

func TestAuthMiddleware_ValidKey(t *testing.T) {
	srv := New(8080, "test-key", "output")
	called := false
	handler := srv.authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer test-key")
	w := httptest.NewRecorder()
	handler(w, req)

	if !called {
		t.Fatal("handler should be called with valid API key")
	}
}

func TestAuthMiddleware_InvalidKey(t *testing.T) {
	srv := New(8080, "test-key", "output")
	handler := srv.authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called with invalid API key")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer wrong-key")
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestCreateScan_MissingDomain(t *testing.T) {
	srv := New(8080, "", "output")
	body := `{}`
	req := httptest.NewRequest("POST", "/api/scan", strings.NewReader(body))
	w := httptest.NewRecorder()

	srv.handleCreateScan(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestStore_CRUD(t *testing.T) {
	store := NewStore()

	job := &ScanJob{ID: "test-1", Status: StatusQueued, Domain: "example.com"}
	store.Create(job)

	got := store.Get("test-1")
	if got == nil {
		t.Fatal("expected job, got nil")
	}
	if got.Domain != "example.com" {
		t.Fatalf("expected example.com, got %s", got.Domain)
	}

	items := store.List()
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}

	if err := store.Delete("test-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got = store.Get("test-1")
	if got.Status != StatusCancelled {
		t.Fatalf("expected cancelled, got %s", got.Status)
	}

	if err := store.Delete("nonexistent"); err == nil {
		t.Fatal("expected error for nonexistent job")
	}
}

func TestGetScan_NotFound(t *testing.T) {
	srv := New(8080, "", "output")
	req := httptest.NewRequest("GET", "/api/scan/nonexistent", nil)
	req.SetPathValue("id", "nonexistent")
	w := httptest.NewRecorder()

	srv.handleGetScan(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestListScans_Empty(t *testing.T) {
	srv := New(8080, "", "output")
	req := httptest.NewRequest("GET", "/api/scans", nil)
	w := httptest.NewRecorder()

	srv.handleListScans(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var items []ScanListItem
	if err := json.NewDecoder(w.Body).Decode(&items); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("expected 0 items, got %d", len(items))
	}
}
