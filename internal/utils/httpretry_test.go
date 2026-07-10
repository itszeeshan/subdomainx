package utils

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestDoWithRetry_SuccessOnFirstTry(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	req, _ := http.NewRequest("GET", srv.URL, nil)
	resp, err := DoWithRetry(srv.Client(), req, 3, 30)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestDoWithRetry_429ThenSuccess(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if calls.Add(1) == 1 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	req, _ := http.NewRequest("GET", srv.URL, nil)
	resp, err := DoWithRetry(srv.Client(), req, 3, 30)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if calls.Load() != 2 {
		t.Fatalf("expected 2 calls, got %d", calls.Load())
	}
}

func TestDoWithRetry_ExhaustedRetries(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.Header().Set("Retry-After", "0")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer srv.Close()

	req, _ := http.NewRequest("GET", srv.URL, nil)
	_, err := DoWithRetry(srv.Client(), req, 2, 1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	rateLimitErr, ok := err.(*RateLimitError)
	if !ok {
		t.Fatalf("expected *RateLimitError, got %T: %v", err, err)
	}
	if rateLimitErr.StatusCode != 429 {
		t.Fatalf("expected status 429, got %d", rateLimitErr.StatusCode)
	}
	// maxRetries=2 means 3 total attempts (0, 1, 2)
	if calls.Load() != 3 {
		t.Fatalf("expected 3 calls, got %d", calls.Load())
	}
}

func TestDoWithRetry_Non429Passthrough(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("forbidden"))
	}))
	defer srv.Close()

	req, _ := http.NewRequest("GET", srv.URL, nil)
	resp, err := DoWithRetry(srv.Client(), req, 3, 30)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 403 {
		t.Fatalf("expected 403, got %d", resp.StatusCode)
	}
}

func TestDoWithRetry_POSTBodyReplay(t *testing.T) {
	var calls atomic.Int32
	var lastBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		lastBody = string(body)
		if calls.Add(1) == 1 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	payload := `{"query":"test"}`
	req, _ := http.NewRequest("POST", srv.URL, bytes.NewReader([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := DoWithRetry(srv.Client(), req, 3, 30)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if lastBody != payload {
		t.Fatalf("expected body %q on retry, got %q", payload, lastBody)
	}
}

func TestDoWithRetry_RetryAfterSeconds(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if calls.Add(1) == 1 {
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	req, _ := http.NewRequest("GET", srv.URL, nil)
	start := time.Now()
	resp, err := DoWithRetry(srv.Client(), req, 3, 30)
	elapsed := time.Since(start)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// With Retry-After: 1 and ±20% jitter, should wait ~0.8-1.2 seconds
	if elapsed < 700*time.Millisecond {
		t.Fatalf("expected wait ~1s, got %v", elapsed)
	}
}

func TestParseRetryAfter_Empty(t *testing.T) {
	d := ParseRetryAfter("", 0)
	if d != 1*time.Second {
		t.Fatalf("expected 1s for attempt 0, got %v", d)
	}
	d = ParseRetryAfter("", 2)
	if d != 4*time.Second {
		t.Fatalf("expected 4s for attempt 2, got %v", d)
	}
}

func TestParseRetryAfter_Integer(t *testing.T) {
	d := ParseRetryAfter("5", 0)
	if d != 5*time.Second {
		t.Fatalf("expected 5s, got %v", d)
	}
}

func TestParseRetryAfter_InvalidFallback(t *testing.T) {
	d := ParseRetryAfter("garbage", 1)
	if d != 2*time.Second {
		t.Fatalf("expected 2s fallback for attempt 1, got %v", d)
	}
}
