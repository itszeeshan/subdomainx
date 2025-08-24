package enumerator

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/itszeeshan/subdomainx/internal/config"
)

type FierceEnumerator struct{}

func (f *FierceEnumerator) Name() string {
	return "fierce"
}

func (f *FierceEnumerator) Enumerate(ctx context.Context, domain string, cfg *config.Config, ) ([]string, error) {
	// Build fierce command
	args := []string{"--domain", domain, "--output", "/dev/stdout"}

	cmd := exec.CommandContext(ctx, "fierce", args...)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("fierce execution failed: %v", err)
	}

	// Parse output (one subdomain per line)
	lines := strings.Split(string(output), "\n")
	var subdomains []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			subdomains = append(subdomains, line)
			// Cache the DNS result
		}
	}

	return subdomains, nil
}

func init() {
	RegisterEnumerator(&FierceEnumerator{})
}
