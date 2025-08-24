package enumerator

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/itszeeshan/subdomainx/internal/config"
)

type URLScanEnumerator struct {
	apiKey string
	client *http.Client
}

type URLScanResponse struct {
	Results []struct {
		Page struct {
			Domain string `json:"domain"`
		} `json:"page"`
	} `json:"results"`
}

func (u *URLScanEnumerator) Name() string {
	return "urlscan"
}

func (u *URLScanEnumerator) Enumerate(ctx context.Context, domain string, cfg *config.Config, ) ([]string, error) {
	// Build API URL
	url := fmt.Sprintf("https://urlscan.io/api/v1/search/?q=domain:%s", domain)

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Add headers
	req.Header.Set("User-Agent", "SubdomainX/1.0")
	req.Header.Set("Accept", "application/json")

	// Add API key if available
	if u.apiKey != "" {
		req.Header.Set("API-Key", u.apiKey)
	}

	// Make request
	resp, err := u.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("urlscan API request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("urlscan API error: %s - %s", resp.Status, string(body))
	}

	// Parse response
	var apiResp URLScanResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse urlscan response: %v", err)
	}

	// Extract subdomains
	subdomainSet := make(map[string]bool)
	for _, result := range apiResp.Results {
		if result.Page.Domain != "" && strings.HasSuffix(result.Page.Domain, "."+domain) {
			subdomainSet[result.Page.Domain] = true
		}
	}

	// Convert set to slice
	var subdomains []string
	for subdomain := range subdomainSet {
		subdomains = append(subdomains, subdomain)
	}

	return subdomains, nil
}

func init() {
	// Create URLScan enumerator with default settings
	enumerator := &URLScanEnumerator{
		apiKey: getURLScanAPIKey(),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
	RegisterEnumerator(enumerator)
}

// getURLScanAPIKey retrieves the API key from environment variable
func getURLScanAPIKey() string {
	// Try to get from environment variable
	apiKey := strings.TrimSpace(os.Getenv("URLSCAN_API_KEY"))
	if apiKey == "" {
		// Could also try to get from config file in the future
		return ""
	}
	return apiKey
}
