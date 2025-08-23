package tests

import (
	"os"
	"testing"

	"github.com/itszeeshan/subdomainx/internal/utils"
)

func TestFileUtils(t *testing.T) {
	// Test file existence
	if utils.FileExists("nonexistent.txt") {
		t.Error("FileExists should return false for non-existent file")
	}

	// Test directory creation
	testDir := "test_output"
	err := utils.EnsureDirectory(testDir)
	if err != nil {
		t.Errorf("Failed to create directory: %v", err)
	}

	// Clean up
	defer os.RemoveAll(testDir)

	// Test writing and reading lines
	testLines := []string{"line1", "line2", "line3"}
	testFile := "test_file.txt"

	err = utils.WriteLines(testFile, testLines)
	if err != nil {
		t.Errorf("Failed to write lines: %v", err)
	}

	defer os.Remove(testFile)

	readLines, err := utils.ReadLines(testFile)
	if err != nil {
		t.Errorf("Failed to read lines: %v", err)
	}

	if len(readLines) != len(testLines) {
		t.Errorf("Expected %d lines, got %d", len(testLines), len(readLines))
	}

	for i, line := range testLines {
		if readLines[i] != line {
			t.Errorf("Expected line %s, got %s", line, readLines[i])
		}
	}
}

func TestValidation(t *testing.T) {
	// Test domain validation
	validDomains := []string{"example.com", "test.example.org", "sub.domain.co.uk"}
	invalidDomains := []string{"", "invalid", "test.", ".com"}

	for _, domain := range validDomains {
		if err := utils.ValidateDomain(domain); err != nil {
			t.Errorf("Domain %s should be valid: %v", domain, err)
		}
	}

	for _, domain := range invalidDomains {
		if err := utils.ValidateDomain(domain); err == nil {
			t.Errorf("Domain %s should be invalid", domain)
		}
	}

	// Test IP validation
	validIPs := []string{"192.168.1.1", "10.0.0.1", "172.16.0.1"}
	invalidIPs := []string{"", "invalid", "256.256.256.256", "192.168.1"}

	for _, ip := range validIPs {
		if err := utils.ValidateIP(ip); err != nil {
			t.Errorf("IP %s should be valid: %v", ip, err)
		}
	}

	for _, ip := range invalidIPs {
		if err := utils.ValidateIP(ip); err == nil {
			t.Errorf("IP %s should be invalid", ip)
		}
	}

	// Test port validation
	validPorts := []int{1, 80, 443, 8080, 65535}
	invalidPorts := []int{0, -1, 65536}

	for _, port := range validPorts {
		if err := utils.ValidatePort(port); err != nil {
			t.Errorf("Port %d should be valid: %v", port, err)
		}
	}

	for _, port := range invalidPorts {
		if err := utils.ValidatePort(port); err == nil {
			t.Errorf("Port %d should be invalid", port)
		}
	}

	// Test URL validation
	validURLs := []string{"http://example.com", "https://test.org"}
	invalidURLs := []string{"", "invalid", "ftp://example.com"}

	for _, url := range validURLs {
		if err := utils.ValidateURL(url); err != nil {
			t.Errorf("URL %s should be valid: %v", url, err)
		}
	}

	for _, url := range invalidURLs {
		if err := utils.ValidateURL(url); err == nil {
			t.Errorf("URL %s should be invalid", url)
		}
	}
}
