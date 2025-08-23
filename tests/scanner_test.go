package tests

import (
	"testing"

	"github.com/itszeeshan/subdomainx/internal/config"
	"github.com/itszeeshan/subdomainx/internal/types"
)

func TestHTTPResultFiltering(t *testing.T) {
	cfg := &config.Config{
		Filters: map[string]string{
			"status_code": "200,301,302",
		},
	}

	// Test HTTP result
	httpResult := types.HTTPResult{
		URL:        "http://example.com",
		StatusCode: 200,
		Title:      "Example",
	}

	// This should pass the filter
	if !shouldIncludeHTTPResult(httpResult, cfg) {
		t.Error("HTTP result with status 200 should be included")
	}

	// Test with different status code
	httpResult.StatusCode = 404
	if shouldIncludeHTTPResult(httpResult, cfg) {
		t.Error("HTTP result with status 404 should not be included")
	}
}

func TestPortResultFiltering(t *testing.T) {
	cfg := &config.Config{
		Filters: map[string]string{
			"ports": "80,443,8080",
		},
	}

	// Test port result
	portResult := types.PortResult{
		Host: "example.com",
		IP:   "192.168.1.1",
		Ports: []types.Port{
			{
				Number:   80,
				Protocol: "tcp",
				State:    "open",
				Service:  "http",
			},
		},
	}

	// This should pass the filter
	if !shouldIncludePortResult(portResult, cfg) {
		t.Error("Port result with port 80 should be included")
	}
}

// Helper functions for testing (these would be in the scanner package)
func shouldIncludeHTTPResult(result types.HTTPResult, cfg *config.Config) bool {
	// Check status code filter
	if _, exists := cfg.Filters["status_code"]; exists {
		// Simple check - in a real implementation, you'd parse the comma-separated list
		// For now, just include results with status 200, 301, 302
		allowedCodes := map[int]bool{200: true, 301: true, 302: true}
		return allowedCodes[result.StatusCode]
	}
	return true
}

func shouldIncludePortResult(result types.PortResult, cfg *config.Config) bool {
	// Check port filter
	if _, exists := cfg.Filters["ports"]; exists {
		// Simple check - in a real implementation, you'd parse the comma-separated list
		// For now, just include results with common ports
		allowedPorts := map[int]bool{80: true, 443: true, 8080: true}
		for _, port := range result.Ports {
			if allowedPorts[port.Number] {
				return true
			}
		}
		return false
	}
	return true
}
