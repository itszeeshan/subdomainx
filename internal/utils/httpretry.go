package utils

import (
	"bytes"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"
)

// RateLimitError is returned when retries are exhausted on 429 responses.
type RateLimitError struct {
	StatusCode int
	RetryAfter time.Duration
	Message    string
}

func (e *RateLimitError) Error() string {
	return e.Message
}

// DoWithRetry executes an HTTP request, handling 429 Too Many Requests
// responses with exponential backoff, jitter, and Retry-After header parsing.
// It retries up to maxRetries times. For requests with a body (POST), the
// body is buffered and re-created on each retry. Non-429 responses pass
// through untouched for the caller to check StatusOK.
func DoWithRetry(client *http.Client, req *http.Request, maxRetries int, timeout int) (*http.Response, error) {
	// Buffer the request body if present so we can replay it on retries.
	var bodyBytes []byte
	if req.Body != nil {
		var err error
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read request body: %w", err)
		}
		_ = req.Body.Close()
		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}

	maxWait := time.Duration(timeout) * time.Second

	for attempt := range maxRetries + 1 {
		if attempt > 0 && bodyBytes != nil {
			req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusTooManyRequests {
			return resp, nil
		}

		// 429 — close body before retry
		_ = resp.Body.Close()

		if attempt == maxRetries {
			return nil, &RateLimitError{
				StatusCode: 429,
				Message:    fmt.Sprintf("rate limited after %d attempts", maxRetries+1),
			}
		}

		wait := ParseRetryAfter(resp.Header.Get("Retry-After"), attempt)
		if wait > maxWait {
			wait = maxWait
		}

		// Add jitter: ±20%
		jitter := time.Duration(float64(wait) * (0.8 + 0.4*rand.Float64()))
		time.Sleep(jitter)
	}

	return nil, &RateLimitError{StatusCode: 429, Message: "rate limited: retries exhausted"}
}

// ParseRetryAfter parses the Retry-After header value. Supports integer
// seconds and HTTP-date formats (RFC1123, RFC850). Falls back to exponential
// backoff (2^attempt seconds) if the header is missing or unparseable.
func ParseRetryAfter(header string, attempt int) time.Duration {
	if header == "" {
		return time.Duration(1<<uint(attempt)) * time.Second
	}

	if seconds, err := strconv.Atoi(header); err == nil && seconds > 0 {
		return time.Duration(seconds) * time.Second
	}

	for _, layout := range []string{time.RFC1123, time.RFC850, "Mon Jan _2 15:04:05 2006"} {
		if t, err := time.Parse(layout, header); err == nil {
			if d := time.Until(t); d > 0 {
				return d
			}
		}
	}

	return time.Duration(1<<uint(attempt)) * time.Second
}
