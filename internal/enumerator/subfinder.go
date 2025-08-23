package enumerator

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/itszeeshan/subdomainx/internal/cache"
	"github.com/itszeeshan/subdomainx/internal/config"
)

type SubfinderEnumerator struct{}

func (s *SubfinderEnumerator) Name() string {
	return "subfinder"
}

func (s *SubfinderEnumerator) Enumerate(ctx context.Context, domain string, cfg *config.Config, cache *cache.DNSCache) ([]string, error) {
	// Build subfinder command
	args := []string{"-d", domain, "-silent"}
	if cfg.Wordlist != "" {
		args = append(args, "-w", cfg.Wordlist)
	}

	cmd := exec.CommandContext(ctx, "subfinder", args...)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("subfinder execution failed: %v", err)
	}

	// Parse output (one subdomain per line)
	lines := strings.Split(string(output), "\n")
	var subdomains []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			subdomains = append(subdomains, line)
			// DNS resolution will be done later in the main enumerator
		}
	}

	return subdomains, nil
}

func init() {
	RegisterEnumerator(&SubfinderEnumerator{})
}
