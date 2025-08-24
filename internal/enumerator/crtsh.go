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

type CrtShEnumerator struct {
	client *http.Client
}

type CrtShResponse []struct {
	NameValue string `json:"name_value"`
}

func (c *CrtShEnumerator) Name() string {
	return "crtsh"
}

func (c *CrtShEnumerator) Enumerate(ctx context.Context, domain string, cfg *config.Config) ([]string, error) {
	// Build API URL - crt.sh provides a public API
	url := fmt.Sprintf("https://crt.sh/?q=%%25.%s&output=json", domain)

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Add headers
	req.Header.Set("Accept", "application/json")

	// Make request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("crtsh API request failed (check network/proxy): %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("crtsh API error: %s - %s", resp.Status, string(body))
	}

	// Parse response
	var apiResp CrtShResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse crtsh response: %v", err)
	}

	// Extract and deduplicate subdomains
	subdomainSet := make(map[string]bool)
	for _, record := range apiResp {
		if record.NameValue != "" {
			// Split by newlines and wildcards
			names := strings.Split(record.NameValue, "\n")
			for _, name := range names {
				name = strings.TrimSpace(name)
				if name != "" && strings.Contains(name, domain) {
					// Remove wildcards and clean up
					name = strings.ReplaceAll(name, "*.", "")
					name = strings.ReplaceAll(name, "*", "")
					if strings.HasSuffix(name, "."+domain) {
						subdomainSet[name] = true
					}
				}
			}
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
	// Create crt.sh enumerator with default settings
	enumerator := &CrtShEnumerator{
		client: &http.Client{
			Timeout: 120 * time.Second, // Increased timeout for API calls
		},
	}
	RegisterEnumerator(enumerator)
}
