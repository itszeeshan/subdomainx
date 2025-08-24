package enumerator

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/itszeeshan/subdomainx/internal/config"
)

type KnockpyEnumerator struct{}

func (k *KnockpyEnumerator) Name() string {
	return "knockpy"
}

func (k *KnockpyEnumerator) Enumerate(ctx context.Context, domain string, cfg *config.Config, ) ([]string, error) {
	// Build knockpy command
	args := []string{domain, "--no-http"}

	cmd := exec.CommandContext(ctx, "knockpy", args...)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("knockpy execution failed: %v", err)
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
	RegisterEnumerator(&KnockpyEnumerator{})
}
