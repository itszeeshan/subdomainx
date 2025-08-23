package enumerator

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/itszeeshan/subdomainx/internal/cache"
	"github.com/itszeeshan/subdomainx/internal/config"
)

type Sublist3rEnumerator struct{}

func (s *Sublist3rEnumerator) Name() string {
	return "sublist3r"
}

func (s *Sublist3rEnumerator) Enumerate(ctx context.Context, domain string, cfg *config.Config, cache *cache.DNSCache) ([]string, error) {
	// Build sublist3r command
	args := []string{"-d", domain, "-o", "/dev/stdout"}

	cmd := exec.CommandContext(ctx, "sublist3r", args...)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("sublist3r execution failed: %v", err)
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
	RegisterEnumerator(&Sublist3rEnumerator{})
}
