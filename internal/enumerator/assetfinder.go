package enumerator

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/itszeeshan/subdomainx/internal/cache"
	"github.com/itszeeshan/subdomainx/internal/config"
)

type AssetfinderEnumerator struct{}

func (a *AssetfinderEnumerator) Name() string {
	return "assetfinder"
}

func (a *AssetfinderEnumerator) Enumerate(ctx context.Context, domain string, cfg *config.Config, cache *cache.DNSCache) ([]string, error) {
	// Build assetfinder command
	args := []string{domain}

	cmd := exec.CommandContext(ctx, "assetfinder", args...)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("assetfinder execution failed: %v", err)
	}

	// Parse output (one subdomain per line)
	lines := strings.Split(string(output), "\n")
	var subdomains []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			subdomains = append(subdomains, line)
			// Cache the DNS result
			cache.Store(line, []string{}) // IPs would be populated by massdns later
		}
	}

	return subdomains, nil
}

func init() {
	RegisterEnumerator(&AssetfinderEnumerator{})
}
