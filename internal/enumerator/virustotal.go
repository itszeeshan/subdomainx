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

	"github.com/itszeeshan/subdomainx/internal/cache"
	"github.com/itszeeshan/subdomainx/internal/config"
)

type VirusTotalEnumerator struct {
	apiKey string
	client *http.Client
}

type VirusTotalResponse struct {
	Data []struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"data"`
	Meta struct {
		Count int `json:"count"`
	} `json:"meta"`
}

func (v *VirusTotalEnumerator) Name() string {
	return "virustotal"
}

func (v *VirusTotalEnumerator) Enumerate(ctx context.Context, domain string, cfg *config.Config, cache *cache.DNSCache) ([]string, error) {
	// Check if API key is configured
	if v.apiKey == "" {
		return nil, fmt.Errorf("virustotal API key not configured")
	}

	// Build API URL
	url := fmt.Sprintf("https://www.virustotal.com/api/v3/domains/%s/subdomains", domain)

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Add headers
	req.Header.Set("x-apikey", v.apiKey)
	req.Header.Set("Accept", "application/json")

	// Make request
	resp, err := v.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("virustotal API request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("virustotal API error: %s - %s", resp.Status, string(body))
	}

	// Parse response
	var apiResp VirusTotalResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse virustotal response: %v", err)
	}

	// Extract subdomains from response
	var subdomains []string
	for _, item := range apiResp.Data {
		if item.Type == "domain" && item.ID != "" {
			// VirusTotal returns full domain names, so we can use them directly
			subdomains = append(subdomains, item.ID)
		}
	}

	return subdomains, nil
}

func init() {
	// Create VirusTotal enumerator with default settings
	enumerator := &VirusTotalEnumerator{
		apiKey: getVirusTotalAPIKey(),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
	RegisterEnumerator(enumerator)
}

// getVirusTotalAPIKey retrieves the API key from environment variable
func getVirusTotalAPIKey() string {
	// Try to get from environment variable
	apiKey := strings.TrimSpace(os.Getenv("VIRUSTOTAL_API_KEY"))
	if apiKey == "" {
		return ""
	}
	return apiKey
}
