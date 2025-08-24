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

type SecurityTrailsEnumerator struct {
	apiKey string
	client *http.Client
}

type SecurityTrailsResponse struct {
	Subdomains []string `json:"subdomains"`
	Meta       struct {
		Total int `json:"total"`
	} `json:"meta"`
}

func (s *SecurityTrailsEnumerator) Name() string {
	return "securitytrails"
}

func (s *SecurityTrailsEnumerator) Enumerate(ctx context.Context, domain string, cfg *config.Config, ) ([]string, error) {
	// Check if API key is configured
	if s.apiKey == "" {
		return nil, fmt.Errorf("securitytrails API key not configured")
	}

	// Build API URL
	url := fmt.Sprintf("https://api.securitytrails.com/v1/domain/%s/subdomains", domain)

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Add headers
	req.Header.Set("APIKEY", s.apiKey)
	req.Header.Set("Accept", "application/json")

	// Make request
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("securitytrails API request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("securitytrails API error: %s - %s", resp.Status, string(body))
	}

	// Parse response
	var apiResp SecurityTrailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse securitytrails response: %v", err)
	}

	// Convert subdomains to full domain names
	var subdomains []string
	for _, subdomain := range apiResp.Subdomains {
		if subdomain != "" {
			fullDomain := fmt.Sprintf("%s.%s", subdomain, domain)
			subdomains = append(subdomains, fullDomain)
		}
	}

	return subdomains, nil
}

func init() {
	// Create SecurityTrails enumerator with default settings
	// API key should be configured via environment variable or config file
	enumerator := &SecurityTrailsEnumerator{
		apiKey: getSecurityTrailsAPIKey(),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
	RegisterEnumerator(enumerator)
}

// getSecurityTrailsAPIKey retrieves the API key from environment variable
func getSecurityTrailsAPIKey() string {
	// Try to get from environment variable
	apiKey := strings.TrimSpace(os.Getenv("SECURITYTRAILS_API_KEY"))
	if apiKey == "" {
		// Could also try to get from config file in the future
		return ""
	}
	return apiKey
}
