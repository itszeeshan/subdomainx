package utils

import (
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/itszeeshan/subdomainx/internal/config"
)

// ValidateInput validates the configuration and input parameters
func ValidateInput(cfg *config.Config) error {
	// Validate wildcard file exists
	if !FileExists(cfg.WildcardFile) {
		return fmt.Errorf("wildcard file not found: %s", cfg.WildcardFile)
	}

	// Validate output directory can be created
	if err := EnsureDirectory(cfg.OutputDir); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// Validate output format
	validFormats := map[string]bool{
		"json": true,
		"txt":  true,
		"html": true,
	}
	if !validFormats[cfg.OutputFormat] {
		return fmt.Errorf("invalid output format: %s. Supported formats: json, txt, html", cfg.OutputFormat)
	}

	// Validate threads
	if cfg.Threads <= 0 {
		return fmt.Errorf("threads must be greater than 0")
	}

	// Validate retries
	if cfg.Retries < 0 {
		return fmt.Errorf("retries cannot be negative")
	}

	// Validate timeout
	if cfg.Timeout <= 0 {
		return fmt.Errorf("timeout must be greater than 0")
	}

	// Validate rate limit
	if cfg.RateLimit <= 0 {
		return fmt.Errorf("rate limit must be greater than 0")
	}

	// Validate wordlist if specified
	if cfg.Wordlist != "" && !FileExists(cfg.Wordlist) {
		return fmt.Errorf("wordlist file not found: %s", cfg.Wordlist)
	}

	return nil
}

// ValidateDomain validates if a string is a valid domain name
func ValidateDomain(domain string) error {
	if domain == "" {
		return fmt.Errorf("domain cannot be empty")
	}

	// Basic domain validation regex
	domainRegex := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*$`)
	if !domainRegex.MatchString(domain) {
		return fmt.Errorf("invalid domain format: %s", domain)
	}

	// Check for valid TLD
	parts := strings.Split(domain, ".")
	if len(parts) < 2 {
		return fmt.Errorf("domain must have at least one subdomain and TLD: %s", domain)
	}

	return nil
}

// ValidateIP validates if a string is a valid IP address
func ValidateIP(ip string) error {
	if net.ParseIP(ip) == nil {
		return fmt.Errorf("invalid IP address: %s", ip)
	}
	return nil
}

// ValidatePort validates if a port number is valid
func ValidatePort(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535, got: %d", port)
	}
	return nil
}

// ValidateURL validates if a string is a valid URL
func ValidateURL(url string) error {
	if url == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	// Basic URL validation
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return fmt.Errorf("URL must start with http:// or https://: %s", url)
	}

	return nil
}
