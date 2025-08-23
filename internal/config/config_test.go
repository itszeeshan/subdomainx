package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigDefault(t *testing.T) {
	// Test loading default config (should work even without config file)
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg == nil {
		t.Fatal("Expected config, got nil")
	}

	// Check default values
	if cfg.OutputDir != "output" {
		t.Errorf("Expected OutputDir 'output', got '%s'", cfg.OutputDir)
	}
	if cfg.OutputFormat != "json" {
		t.Errorf("Expected OutputFormat 'json', got '%s'", cfg.OutputFormat)
	}
	if cfg.Threads != 10 {
		t.Errorf("Expected Threads 10, got %d", cfg.Threads)
	}
	if cfg.Retries != 3 {
		t.Errorf("Expected Retries 3, got %d", cfg.Retries)
	}
	if cfg.Timeout != 30 {
		t.Errorf("Expected Timeout 30, got %d", cfg.Timeout)
	}
	if cfg.RateLimit != 100 {
		t.Errorf("Expected RateLimit 100, got %d", cfg.RateLimit)
	}
	if cfg.Tools == nil {
		t.Error("Expected Tools map to be initialized")
	}
	if cfg.Filters == nil {
		t.Error("Expected Filters map to be initialized")
	}
}

func TestLoadConfigFromFile(t *testing.T) {
	// Create a temporary config file
	tmpDir, err := os.MkdirTemp("", "test_config")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "test_config.yaml")
	configContent := `output_dir: custom_output
output_format: html
threads: 20
retries: 5
timeout: 60
rate_limit: 200
tools:
  subfinder: true
  amass: false
filters:
  status_codes: "200,301,302"
  ports: "80,443,8080"`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Load config from file
	cfg, err := LoadConfigFromFile(configPath)
	if err != nil {
		t.Fatalf("LoadConfigFromFile failed: %v", err)
	}

	if cfg == nil {
		t.Fatal("Expected config, got nil")
	}

	// Check loaded values
	if cfg.OutputDir != "custom_output" {
		t.Errorf("Expected OutputDir 'custom_output', got '%s'", cfg.OutputDir)
	}
	if cfg.OutputFormat != "html" {
		t.Errorf("Expected OutputFormat 'html', got '%s'", cfg.OutputFormat)
	}
	if cfg.Threads != 20 {
		t.Errorf("Expected Threads 20, got %d", cfg.Threads)
	}
	if cfg.Retries != 5 {
		t.Errorf("Expected Retries 5, got %d", cfg.Retries)
	}
	if cfg.Timeout != 60 {
		t.Errorf("Expected Timeout 60, got %d", cfg.Timeout)
	}
	if cfg.RateLimit != 200 {
		t.Errorf("Expected RateLimit 200, got %d", cfg.RateLimit)
	}

	// Check tools map
	if !cfg.Tools["subfinder"] {
		t.Error("Expected subfinder to be true")
	}
	if cfg.Tools["amass"] {
		t.Error("Expected amass to be false")
	}

	// Check filters map
	if cfg.Filters["status_codes"] != "200,301,302" {
		t.Errorf("Expected status_codes '200,301,302', got '%s'", cfg.Filters["status_codes"])
	}
	if cfg.Filters["ports"] != "80,443,8080" {
		t.Errorf("Expected ports '80,443,8080', got '%s'", cfg.Filters["ports"])
	}
}

func TestLoadConfigFromFileNonExistent(t *testing.T) {
	// Test loading from non-existent file (should use defaults)
	cfg, err := LoadConfigFromFile("non_existent_config.yaml")
	if err != nil {
		t.Fatalf("LoadConfigFromFile failed: %v", err)
	}

	if cfg == nil {
		t.Fatal("Expected config, got nil")
	}

	// Should have default values
	if cfg.OutputDir != "output" {
		t.Errorf("Expected default OutputDir 'output', got '%s'", cfg.OutputDir)
	}
	if cfg.OutputFormat != "json" {
		t.Errorf("Expected default OutputFormat 'json', got '%s'", cfg.OutputFormat)
	}
}

func TestLoadConfigFromFileInvalidYAML(t *testing.T) {
	// Create a temporary config file with invalid YAML
	tmpDir, err := os.MkdirTemp("", "test_config")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "invalid_config.yaml")
	invalidYAML := `output_dir: "unclosed_quote
threads: invalid_number
output_format: html`

	if err := os.WriteFile(configPath, []byte(invalidYAML), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Load config from invalid file
	_, err = LoadConfigFromFile(configPath)
	if err == nil {
		t.Error("Expected error for invalid YAML, got nil")
	}

	expectedMsg := "failed to parse config file:"
	if err.Error()[:len(expectedMsg)] != expectedMsg {
		t.Errorf("Expected error message starting with '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestLoadConfigFromFileReadError(t *testing.T) {
	// Test with a directory path (should cause read error)
	tmpDir, err := os.MkdirTemp("", "test_config")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Try to load a directory as a config file
	_, err = LoadConfigFromFile(tmpDir)
	if err == nil {
		t.Error("Expected error for directory path, got nil")
	}

	expectedMsg := "failed to read config file:"
	if err.Error()[:len(expectedMsg)] != expectedMsg {
		t.Errorf("Expected error message starting with '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestConfigSave(t *testing.T) {
	// Create a temporary directory for config
	tmpDir, err := os.MkdirTemp("", "test_config")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create configs subdirectory
	configsDir := filepath.Join(tmpDir, "configs")
	if err := os.MkdirAll(configsDir, 0755); err != nil {
		t.Fatalf("Failed to create configs directory: %v", err)
	}

	// Change to temp directory temporarily
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Create config with custom values
	cfg := &Config{
		OutputDir:    "custom_output",
		OutputFormat: "html",
		Threads:      25,
		Retries:      7,
		Timeout:      90,
		RateLimit:    300,
		Tools: map[string]bool{
			"subfinder": true,
			"amass":     false,
		},
		Filters: map[string]string{
			"status_codes": "200,301,302,404",
			"ports":        "80,443,8080,8443",
		},
	}

	// Save config
	err = cfg.Save()
	if err != nil {
		t.Fatalf("Config.Save failed: %v", err)
	}

	// Verify file was created
	configFile := filepath.Join("configs", "default.yaml")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}

	// Load the saved config to verify
	loadedCfg, err := LoadConfigFromFile(configFile)
	if err != nil {
		t.Fatalf("Failed to load saved config: %v", err)
	}

	// Verify values were saved correctly
	if loadedCfg.OutputDir != cfg.OutputDir {
		t.Errorf("Expected OutputDir '%s', got '%s'", cfg.OutputDir, loadedCfg.OutputDir)
	}
	if loadedCfg.OutputFormat != cfg.OutputFormat {
		t.Errorf("Expected OutputFormat '%s', got '%s'", cfg.OutputFormat, loadedCfg.OutputFormat)
	}
	if loadedCfg.Threads != cfg.Threads {
		t.Errorf("Expected Threads %d, got %d", cfg.Threads, loadedCfg.Threads)
	}
	if loadedCfg.Retries != cfg.Retries {
		t.Errorf("Expected Retries %d, got %d", cfg.Retries, loadedCfg.Retries)
	}
	if loadedCfg.Timeout != cfg.Timeout {
		t.Errorf("Expected Timeout %d, got %d", cfg.Timeout, loadedCfg.Timeout)
	}
	if loadedCfg.RateLimit != cfg.RateLimit {
		t.Errorf("Expected RateLimit %d, got %d", cfg.RateLimit, loadedCfg.RateLimit)
	}

	// Verify tools map
	if loadedCfg.Tools["subfinder"] != cfg.Tools["subfinder"] {
		t.Errorf("Expected subfinder %v, got %v", cfg.Tools["subfinder"], loadedCfg.Tools["subfinder"])
	}
	if loadedCfg.Tools["amass"] != cfg.Tools["amass"] {
		t.Errorf("Expected amass %v, got %v", cfg.Tools["amass"], loadedCfg.Tools["amass"])
	}

	// Verify filters map
	if loadedCfg.Filters["status_codes"] != cfg.Filters["status_codes"] {
		t.Errorf("Expected status_codes '%s', got '%s'", cfg.Filters["status_codes"], loadedCfg.Filters["status_codes"])
	}
	if loadedCfg.Filters["ports"] != cfg.Filters["ports"] {
		t.Errorf("Expected ports '%s', got '%s'", cfg.Filters["ports"], loadedCfg.Filters["ports"])
	}
}

func TestConfigSaveMarshalError(t *testing.T) {
	// Create a config that would cause marshaling issues
	// This is hard to trigger with simple structs, but we can test the error path
	// by creating a config with a channel (which can't be marshaled to YAML)

	// For now, we'll test the write error path by using a non-writable directory
	tmpDir, err := os.MkdirTemp("", "test_config")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create configs subdirectory with read-only permissions
	configsDir := filepath.Join(tmpDir, "configs")
	if err := os.MkdirAll(configsDir, 0444); err != nil {
		t.Fatalf("Failed to create configs directory: %v", err)
	}

	// Change to temp directory temporarily
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	cfg := &Config{
		OutputDir:    "test",
		OutputFormat: "json",
		Threads:      10,
		Retries:      3,
		Timeout:      30,
		RateLimit:    100,
		Tools:        make(map[string]bool),
		Filters:      make(map[string]string),
	}

	// Try to save config (should fail due to read-only directory)
	err = cfg.Save()
	if err == nil {
		t.Error("Expected error for read-only directory, got nil")
	}

	expectedMsg := "failed to write config file:"
	if err.Error()[:len(expectedMsg)] != expectedMsg {
		t.Errorf("Expected error message starting with '%s', got '%s'", expectedMsg, err.Error())
	}
}
