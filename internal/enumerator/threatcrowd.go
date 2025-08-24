package enumerator

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/itszeeshan/subdomainx/internal/config"
)

type ThreatCrowdEnumerator struct {
	client *http.Client
}

type ThreatCrowdResponse struct {
	Subdomains   []string `json:"subdomains"`
	ResponseCode string   `json:"response_code"`
}

func (t *ThreatCrowdEnumerator) Name() string {
	return "threatcrowd"
}

func (t *ThreatCrowdEnumerator) Enumerate(ctx context.Context, domain string, cfg *config.Config) ([]string, error) {
	// Build API URL - ThreatCrowd provides a public API
	url := fmt.Sprintf("https://www.threatcrowd.org/searchApi/v2/domain/report/?domain=%s", domain)

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Add headers
	req.Header.Set("User-Agent", "SubdomainX/1.0")
	req.Header.Set("Accept", "application/json")

	// Make request
	resp, err := t.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("threatcrowd API request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("threatcrowd API error: %s - %s", resp.Status, string(body))
	}

	// Parse response
	var apiResp ThreatCrowdResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse threatcrowd response: %v", err)
	}

	// Check if response is valid
	if apiResp.ResponseCode != "1" {
		return nil, fmt.Errorf("threatcrowd API returned error response: %s", apiResp.ResponseCode)
	}

	// Extract subdomains
	var subdomains []string
	for _, subdomain := range apiResp.Subdomains {
		if subdomain != "" && strings.HasSuffix(subdomain, "."+domain) {
			subdomains = append(subdomains, subdomain)
		}
	}

	return subdomains, nil
}

func init() {
	// Create ThreatCrowd enumerator with default settings
	enumerator := &ThreatCrowdEnumerator{
		client: &http.Client{
			Timeout: 60 * time.Second, // Increased timeout for API calls
		},
	}
	RegisterEnumerator(enumerator)
}
