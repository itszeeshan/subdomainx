package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/itszeeshan/subdomainx/internal/config"
)

func TestValidateInput(t *testing.T) {
	// Create temporary files and directories for testing
	tmpDir, err := os.MkdirTemp("", "test_validate")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test wildcard file
	wildcardFile := filepath.Join(tmpDir, "domains.txt")
	if err := os.WriteFile(wildcardFile, []byte("example.com\n"), 0644); err != nil {
		t.Fatalf("Failed to create wildcard file: %v", err)
	}

	// Create output directory
	outputDir := filepath.Join(tmpDir, "output")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	// Test valid configuration
	cfg := &config.Config{
		WildcardFile: wildcardFile,
		OutputDir:    outputDir,
		OutputFormat: "json",
		Threads:      10,
		Retries:      3,
		Timeout:      30,
		RateLimit:    100,
		Tools:        make(map[string]bool),
		Filters:      make(map[string]string),
	}

	err = ValidateInput(cfg)
	if err != nil {
		t.Errorf("Expected no error for valid config, got %v", err)
	}
}

func TestValidateInputMissingWildcardFile(t *testing.T) {
	cfg := &config.Config{
		WildcardFile: "non_existent_file.txt",
		OutputDir:    "output",
		OutputFormat: "json",
		Threads:      10,
		Retries:      3,
		Timeout:      30,
		RateLimit:    100,
		Tools:        make(map[string]bool),
		Filters:      make(map[string]string),
	}

	err := ValidateInput(cfg)
	if err == nil {
		t.Error("Expected error for missing wildcard file")
	}

	expectedMsg := "wildcard file not found: non_existent_file.txt"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestValidateInputInvalidOutputFormat(t *testing.T) {
	// Create temporary files
	tmpDir, err := os.MkdirTemp("", "test_validate")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	wildcardFile := filepath.Join(tmpDir, "domains.txt")
	if err := os.WriteFile(wildcardFile, []byte("example.com\n"), 0644); err != nil {
		t.Fatalf("Failed to create wildcard file: %v", err)
	}

	cfg := &config.Config{
		WildcardFile: wildcardFile,
		OutputDir:    "output",
		OutputFormat: "invalid",
		Threads:      10,
		Retries:      3,
		Timeout:      30,
		RateLimit:    100,
		Tools:        make(map[string]bool),
		Filters:      make(map[string]string),
	}

	err = ValidateInput(cfg)
	if err == nil {
		t.Error("Expected error for invalid output format")
	}

	expectedMsg := "invalid output format: invalid. Supported formats: json, txt, html"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestValidateInputInvalidThreads(t *testing.T) {
	// Create temporary files
	tmpDir, err := os.MkdirTemp("", "test_validate")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	wildcardFile := filepath.Join(tmpDir, "domains.txt")
	if err := os.WriteFile(wildcardFile, []byte("example.com\n"), 0644); err != nil {
		t.Fatalf("Failed to create wildcard file: %v", err)
	}

	cfg := &config.Config{
		WildcardFile: wildcardFile,
		OutputDir:    "output",
		OutputFormat: "json",
		Threads:      0, // Invalid
		Retries:      3,
		Timeout:      30,
		RateLimit:    100,
		Tools:        make(map[string]bool),
		Filters:      make(map[string]string),
	}

	err = ValidateInput(cfg)
	if err == nil {
		t.Error("Expected error for invalid threads")
	}

	expectedMsg := "threads must be greater than 0"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestValidateInputInvalidRetries(t *testing.T) {
	// Create temporary files
	tmpDir, err := os.MkdirTemp("", "test_validate")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	wildcardFile := filepath.Join(tmpDir, "domains.txt")
	if err := os.WriteFile(wildcardFile, []byte("example.com\n"), 0644); err != nil {
		t.Fatalf("Failed to create wildcard file: %v", err)
	}

	cfg := &config.Config{
		WildcardFile: wildcardFile,
		OutputDir:    "output",
		OutputFormat: "json",
		Threads:      10,
		Retries:      -1, // Invalid
		Timeout:      30,
		RateLimit:    100,
		Tools:        make(map[string]bool),
		Filters:      make(map[string]string),
	}

	err = ValidateInput(cfg)
	if err == nil {
		t.Error("Expected error for invalid retries")
	}

	expectedMsg := "retries cannot be negative"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestValidateInputInvalidTimeout(t *testing.T) {
	// Create temporary files
	tmpDir, err := os.MkdirTemp("", "test_validate")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	wildcardFile := filepath.Join(tmpDir, "domains.txt")
	if err := os.WriteFile(wildcardFile, []byte("example.com\n"), 0644); err != nil {
		t.Fatalf("Failed to create wildcard file: %v", err)
	}

	cfg := &config.Config{
		WildcardFile: wildcardFile,
		OutputDir:    "output",
		OutputFormat: "json",
		Threads:      10,
		Retries:      3,
		Timeout:      0, // Invalid
		RateLimit:    100,
		Tools:        make(map[string]bool),
		Filters:      make(map[string]string),
	}

	err = ValidateInput(cfg)
	if err == nil {
		t.Error("Expected error for invalid timeout")
	}

	expectedMsg := "timeout must be greater than 0"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestValidateInputInvalidRateLimit(t *testing.T) {
	// Create temporary files
	tmpDir, err := os.MkdirTemp("", "test_validate")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	wildcardFile := filepath.Join(tmpDir, "domains.txt")
	if err := os.WriteFile(wildcardFile, []byte("example.com\n"), 0644); err != nil {
		t.Fatalf("Failed to create wildcard file: %v", err)
	}

	cfg := &config.Config{
		WildcardFile: wildcardFile,
		OutputDir:    "output",
		OutputFormat: "json",
		Threads:      10,
		Retries:      3,
		Timeout:      30,
		RateLimit:    0, // Invalid
		Tools:        make(map[string]bool),
		Filters:      make(map[string]string),
	}

	err = ValidateInput(cfg)
	if err == nil {
		t.Error("Expected error for invalid rate limit")
	}

	expectedMsg := "rate limit must be greater than 0"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestValidateInputMissingWordlist(t *testing.T) {
	// Create temporary files
	tmpDir, err := os.MkdirTemp("", "test_validate")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	wildcardFile := filepath.Join(tmpDir, "domains.txt")
	if err := os.WriteFile(wildcardFile, []byte("example.com\n"), 0644); err != nil {
		t.Fatalf("Failed to create wildcard file: %v", err)
	}

	cfg := &config.Config{
		WildcardFile: wildcardFile,
		OutputDir:    "output",
		OutputFormat: "json",
		Threads:      10,
		Retries:      3,
		Timeout:      30,
		RateLimit:    100,
		Wordlist:     "non_existent_wordlist.txt", // Invalid
		Tools:        make(map[string]bool),
		Filters:      make(map[string]string),
	}

	err = ValidateInput(cfg)
	if err == nil {
		t.Error("Expected error for missing wordlist")
	}

	expectedMsg := "wordlist file not found: non_existent_wordlist.txt"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestValidateDomain(t *testing.T) {
	// Test valid domains
	validDomains := []string{
		"example.com",
		"sub.example.com",
		"test-sub.domain.co.uk",
		"a.b.c.d.example.org",
		"123.example.com",
		"example123.com",
	}

	for _, domain := range validDomains {
		err := ValidateDomain(domain)
		if err != nil {
			t.Errorf("Expected no error for valid domain '%s', got %v", domain, err)
		}
	}
}

func TestValidateDomainInvalid(t *testing.T) {
	// Test invalid domains
	invalidDomains := []string{
		"",             // Empty
		"example",      // No TLD
		".example.com", // Starts with dot
		"example.com.", // Ends with dot
		"example..com", // Double dot
		"-example.com", // Starts with dash
		"example-.com", // Ends with dash
		"example.com-", // Ends with dash
		"example com",  // Space
		"example@com",  // Invalid character
	}

	for _, domain := range invalidDomains {
		err := ValidateDomain(domain)
		if err == nil {
			t.Errorf("Expected error for invalid domain '%s'", domain)
		}
	}
}

func TestValidateIP(t *testing.T) {
	// Test valid IPs
	validIPs := []string{
		"192.168.1.1",
		"10.0.0.1",
		"172.16.0.1",
		"127.0.0.1",
		"::1",
		"2001:db8::1",
		"fe80::1",
	}

	for _, ip := range validIPs {
		err := ValidateIP(ip)
		if err != nil {
			t.Errorf("Expected no error for valid IP '%s', got %v", ip, err)
		}
	}
}

func TestValidateIPInvalid(t *testing.T) {
	// Test invalid IPs
	invalidIPs := []string{
		"",
		"256.1.2.3",
		"1.2.3.256",
		"1.2.3.4.5",
		"1.2.3",
		"invalid",
		"192.168.1.1.1",
	}

	for _, ip := range invalidIPs {
		err := ValidateIP(ip)
		if err == nil {
			t.Errorf("Expected error for invalid IP '%s'", ip)
		}
	}
}

func TestValidatePort(t *testing.T) {
	// Test valid ports
	validPorts := []int{1, 80, 443, 8080, 65535}

	for _, port := range validPorts {
		err := ValidatePort(port)
		if err != nil {
			t.Errorf("Expected no error for valid port %d, got %v", port, err)
		}
	}
}

func TestValidatePortInvalid(t *testing.T) {
	// Test invalid ports
	invalidPorts := []int{0, -1, 65536, 99999}

	for _, port := range invalidPorts {
		err := ValidatePort(port)
		if err == nil {
			t.Errorf("Expected error for invalid port %d", port)
		}
	}
}

func TestValidateURL(t *testing.T) {
	// Test valid URLs
	validURLs := []string{
		"http://example.com",
		"https://example.com",
		"http://sub.example.com/path",
		"https://example.com:8080",
		"http://example.com/path?param=value",
	}

	for _, url := range validURLs {
		err := ValidateURL(url)
		if err != nil {
			t.Errorf("Expected no error for valid URL '%s', got %v", url, err)
		}
	}
}

func TestValidateURLInvalid(t *testing.T) {
	// Test invalid URLs
	invalidURLs := []string{
		"",
		"example.com",
		"ftp://example.com",
		"://example.com",
		"http:/example.com",
		"https:/example.com",
	}

	for _, url := range invalidURLs {
		err := ValidateURL(url)
		if err == nil {
			t.Errorf("Expected error for invalid URL '%s'", url)
		}
	}
}
