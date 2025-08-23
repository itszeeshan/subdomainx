package enumerator

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/itszeeshan/subdomainx/internal/cache"
	"github.com/itszeeshan/subdomainx/internal/config"
)

type AmassEnumerator struct{}

func (a *AmassEnumerator) Name() string {
	return "amass"
}

func (a *AmassEnumerator) Enumerate(ctx context.Context, domain string, cfg *config.Config, cache *cache.DNSCache) ([]string, error) {
	// Build amass command
	args := []string{"enum", "-d", domain, "-json", "/dev/stdout"}
	if cfg.Wordlist != "" {
		args = append(args, "-w", cfg.Wordlist)
	}

	cmd := exec.CommandContext(ctx, "amass", args...)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("amass execution failed: %v", err)
	}

	// Parse JSON output
	lines := strings.Split(string(output), "\n")
	var subdomains []string

	for _, line := range lines {
		if line == "" {
			continue
		}

		var result struct {
			Name string `json:"name"`
		}

		if err := json.Unmarshal([]byte(line), &result); err != nil {
			continue
		}

		if result.Name != "" {
			subdomains = append(subdomains, result.Name)
			// Cache the DNS result
			cache.Store(result.Name, []string{}) // IPs would be populated by massdns later
		}
	}

	return subdomains, nil
}

func init() {
	RegisterEnumerator(&AmassEnumerator{})
}
