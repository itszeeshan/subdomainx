package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	WildcardFile   string            `yaml:"wildcard_file" json:"wildcard_file"`
	UniqueName     string            `yaml:"unique_name" json:"unique_name"`
	OutputDir      string            `yaml:"output_dir" json:"output_dir"`
	OutputFormat   string            `yaml:"output_format" json:"output_format"`
	Tools          map[string]bool   `yaml:"tools" json:"tools"`
	Wordlist       string            `yaml:"wordlist" json:"wordlist"`
	Threads        int               `yaml:"threads" json:"threads"`
	Retries        int               `yaml:"retries" json:"retries"`
	Timeout        int               `yaml:"timeout" json:"timeout"`
	RateLimit      int               `yaml:"rate_limit" json:"rate_limit"`
	Filters        map[string]string `yaml:"filters" json:"filters"`
	MaxHTTPTargets int               `yaml:"max_http_targets" json:"max_http_targets"`
}

func LoadConfig() (*Config, error) {
	return LoadConfigFromFile(filepath.Join("configs", "default.yaml"))
}

func LoadConfigFromFile(configPath string) (*Config, error) {
	// Default configuration
	cfg := &Config{
		OutputDir:    "output",
		OutputFormat: "json",
		Threads:      10,
		Retries:      3,
		Timeout:      30,
		RateLimit:    100,
		Wordlist:     "", // No default wordlist required
		Tools:        make(map[string]bool),
		Filters:      make(map[string]string),
	}

	// Load from config file if exists
	if _, err := os.Stat(configPath); err == nil {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %v", err)
		}

		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %v", err)
		}
	}

	return cfg, nil
}

func (c *Config) Save() error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	configPath := filepath.Join("configs", "default.yaml")
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}
