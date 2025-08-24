package enumerator

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/itszeeshan/subdomainx/internal/config"
)

type HackerTargetEnumerator struct {
	apiKey string
	client *http.Client
}

func (h *HackerTargetEnumerator) Name() string {
	return "hackertarget"
}

func (h *HackerTargetEnumerator) Enumerate(ctx context.Context, domain string, cfg *config.Config, ) ([]string, error) {
	// Build API URL
	url := fmt.Sprintf("https://api.hackertarget.com/hostsearch/?q=%s", domain)

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Add headers
	req.Header.Set("User-Agent", "SubdomainX/1.0")
	req.Header.Set("Accept", "text/plain")

	// Add API key if available
	if h.apiKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", h.apiKey))
	}

	// Make request
	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("hackertarget API request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("hackertarget API error: %s - %s", resp.Status, string(body))
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read hackertarget response: %v", err)
	}

	// Parse response (HackerTarget returns plain text with one subdomain per line)
	responseText := string(body)
	lines := strings.Split(responseText, "\n")

	// Extract subdomains
	var subdomains []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "API") && !strings.HasPrefix(line, "error") {
			// HackerTarget format: subdomain.domain.com,IP
			parts := strings.Split(line, ",")
			if len(parts) > 0 {
				subdomain := strings.TrimSpace(parts[0])
				if subdomain != "" && strings.HasSuffix(subdomain, "."+domain) {
					subdomains = append(subdomains, subdomain)
				}
			}
		}
	}

	return subdomains, nil
}

func init() {
	// Create HackerTarget enumerator with default settings
	enumerator := &HackerTargetEnumerator{
		apiKey: getHackerTargetAPIKey(),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
	RegisterEnumerator(enumerator)
}

// getHackerTargetAPIKey retrieves the API key from environment variable
func getHackerTargetAPIKey() string {
	// Try to get from environment variable
	apiKey := strings.TrimSpace(os.Getenv("HACKERTARGET_API_KEY"))
	if apiKey == "" {
		// Could also try to get from config file in the future
		return ""
	}
	return apiKey
}
