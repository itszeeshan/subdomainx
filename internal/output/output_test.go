package output

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/itszeeshan/subdomainx/internal/config"
	"github.com/itszeeshan/subdomainx/internal/types"
)

func TestGenerate(t *testing.T) {
	// Create temporary directory for output
	tmpDir, err := os.MkdirTemp("", "test_output")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test data
	subdomainResults := []types.SubdomainResult{
		{Subdomain: "test1.example.com", Source: "subfinder", IPs: []string{"192.168.1.1"}},
		{Subdomain: "test2.example.com", Source: "amass", IPs: []string{"10.0.0.1"}},
	}
	httpResults := []types.HTTPResult{
		{URL: "https://test1.example.com", StatusCode: 200, Title: "Test Page"},
		{URL: "https://test2.example.com", StatusCode: 404},
	}
	portResults := []types.PortResult{
		{Host: "test1.example.com", IP: "192.168.1.1", Ports: []types.Port{{Number: 80, Protocol: "tcp", State: "open"}}},
	}

	// Test JSON format
	cfg := &config.Config{
		UniqueName:   "test_scan",
		OutputDir:    tmpDir,
		OutputFormat: "json",
	}

	err = Generate(cfg, subdomainResults, httpResults, portResults)
	if err != nil {
		t.Fatalf("Generate JSON failed: %v", err)
	}

	// Verify JSON files were created
	jsonFile := filepath.Join(tmpDir, "test_scan_results.json")
	if _, err := os.Stat(jsonFile); os.IsNotExist(err) {
		t.Error("JSON results file was not created")
	}

	// Test TXT format
	cfg.OutputFormat = "txt"
	err = Generate(cfg, subdomainResults, httpResults, portResults)
	if err != nil {
		t.Fatalf("Generate TXT failed: %v", err)
	}

	// Verify TXT file was created
	txtFile := filepath.Join(tmpDir, "test_scan_subdomains.txt")
	if _, err := os.Stat(txtFile); os.IsNotExist(err) {
		t.Error("TXT file was not created")
	}

	// Test HTML format
	cfg.OutputFormat = "html"
	err = Generate(cfg, subdomainResults, httpResults, portResults)
	if err != nil {
		t.Fatalf("Generate HTML failed: %v", err)
	}

	// Verify HTML file was created
	htmlFile := filepath.Join(tmpDir, "test_scan_report.html")
	if _, err := os.Stat(htmlFile); os.IsNotExist(err) {
		t.Error("HTML file was not created")
	}
}

func TestGenerateInvalidFormat(t *testing.T) {
	// Create temporary directory for output
	tmpDir, err := os.MkdirTemp("", "test_output")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := &config.Config{
		UniqueName:   "test_scan",
		OutputDir:    tmpDir,
		OutputFormat: "invalid",
	}

	// Test invalid format
	err = Generate(cfg, []types.SubdomainResult{}, []types.HTTPResult{}, []types.PortResult{})
	if err == nil {
		t.Error("Expected error for invalid format")
	}

	expectedMsg := "unsupported output format: invalid"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestGenerateEmptyResults(t *testing.T) {
	// Create temporary directory for output
	tmpDir, err := os.MkdirTemp("", "test_output")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cfg := &config.Config{
		UniqueName:   "empty_scan",
		OutputDir:    tmpDir,
		OutputFormat: "json",
	}

	// Test with empty results
	err = Generate(cfg, []types.SubdomainResult{}, []types.HTTPResult{}, []types.PortResult{})
	if err != nil {
		t.Fatalf("Generate JSON with empty results failed: %v", err)
	}

	// Verify JSON file was created
	jsonFile := filepath.Join(tmpDir, "empty_scan_results.json")
	if _, err := os.Stat(jsonFile); os.IsNotExist(err) {
		t.Error("JSON file was not created for empty results")
	}
}
