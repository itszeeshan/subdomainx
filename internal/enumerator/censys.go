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

type CensysEnumerator struct {
	apiID  string
	secret string
	client *http.Client
}

type CensysResponse struct {
	Result struct {
		Hits []struct {
			Names []string `json:"names"`
		} `json:"hits"`
		Total int `json:"total"`
	} `json:"result"`
}

func (c *CensysEnumerator) Name() string {
	return "censys"
}

func (c *CensysEnumerator) Enumerate(ctx context.Context, domain string, cfg *config.Config, ) ([]string, error) {
	// Check if API credentials are configured
	if c.apiID == "" || c.secret == "" {
		return nil, fmt.Errorf("censys API credentials not configured")
	}

	// Build API URL
	url := "https://search.censys.io/api/v2/hosts/search"

	// Create request body
	query := fmt.Sprintf("names:*.%s", domain)
	requestBody := fmt.Sprintf(`{"q":"%s","per_page":100}`, query)

	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Add headers
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s:%s", c.apiID, c.secret))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Make request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("censys API request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("censys API error: %s - %s", resp.Status, string(body))
	}

	// Parse response
	var apiResp CensysResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse censys response: %v", err)
	}

	// Extract subdomains from response
	subdomainSet := make(map[string]bool)
	for _, hit := range apiResp.Result.Hits {
		for _, name := range hit.Names {
			if strings.HasSuffix(name, "."+domain) && name != domain {
				subdomainSet[name] = true
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
	// Create Censys enumerator with default settings
	enumerator := &CensysEnumerator{
		apiID:  getCensysAPIID(),
		secret: getCensysSecret(),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
	RegisterEnumerator(enumerator)
}

// getCensysAPIID retrieves the API ID from environment variable
func getCensysAPIID() string {
	return strings.TrimSpace(os.Getenv("CENSYS_API_ID"))
}

// getCensysSecret retrieves the API secret from environment variable
func getCensysSecret() string {
	return strings.TrimSpace(os.Getenv("CENSYS_SECRET"))
}
