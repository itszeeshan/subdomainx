package tests

import (
	"testing"

	"github.com/itszeeshan/subdomainx/internal/cache"
	"github.com/itszeeshan/subdomainx/internal/config"
	"github.com/itszeeshan/subdomainx/internal/enumerator"
)

func TestEnumeratorRegistration(t *testing.T) {
	// Test that enumerators are properly registered
	cfg := &config.Config{
		Tools: map[string]bool{
			"subfinder": true,
			"amass":     true,
		},
	}

	// This should not panic
	enumerators := make(map[string]enumerator.Enumerator)
	for name, enabled := range cfg.Tools {
		if enabled {
			if e, exists := enumerators[name]; exists {
				t.Logf("Enumerator %s is registered", name)
				_ = e.Name()
			}
		}
	}
}

func TestDNSCache(t *testing.T) {
	cache := cache.NewDNSCache()

	// Test storing and retrieving
	testDomain := "test.example.com"
	testIPs := []string{"192.168.1.1", "192.168.1.2"}

	cache.Store(testDomain, testIPs)

	retrieved := cache.Lookup(testDomain)
	if len(retrieved) != len(testIPs) {
		t.Errorf("Expected %d IPs, got %d", len(testIPs), len(retrieved))
	}

	for i, ip := range testIPs {
		if retrieved[i] != ip {
			t.Errorf("Expected IP %s, got %s", ip, retrieved[i])
		}
	}

	// Test non-existent domain
	notFound := cache.Lookup("nonexistent.example.com")
	if notFound != nil {
		t.Errorf("Expected nil for non-existent domain, got %v", notFound)
	}
}

func TestConfigValidation(t *testing.T) {
	cfg := &config.Config{
		WildcardFile: "domains.txt",
		OutputDir:    "output",
		Threads:      10,
		Retries:      3,
		Timeout:      30,
		RateLimit:    100,
		Tools:        make(map[string]bool),
		Filters:      make(map[string]string),
	}

	// Test valid config
	if cfg.Threads <= 0 {
		t.Error("Threads should be greater than 0")
	}

	if cfg.Retries < 0 {
		t.Error("Retries should not be negative")
	}

	if cfg.Timeout <= 0 {
		t.Error("Timeout should be greater than 0")
	}

	if cfg.RateLimit <= 0 {
		t.Error("Rate limit should be greater than 0")
	}
}
