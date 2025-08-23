package enumerator

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/itszeeshan/subdomainx/internal/cache"
	"github.com/itszeeshan/subdomainx/internal/config"
)

type AltDNSEnumerator struct{}

func (a *AltDNSEnumerator) Name() string {
	return "altdns"
}

func (a *AltDNSEnumerator) Enumerate(ctx context.Context, domain string, cfg *config.Config, cache *cache.DNSCache) ([]string, error) {
	// Build altdns command
	args := []string{"-i", "/dev/stdin", "-o", "/dev/stdout"}

	// Use custom wordlist if provided, otherwise use default
	if cfg.Wordlist != "" {
		args = append(args, "-w", cfg.Wordlist)
	} else {
		args = append(args, "-w", "/usr/share/altdns/words.txt")
	}

	cmd := exec.CommandContext(ctx, "altdns", args...)

	// Create input with base subdomains
	baseSubdomains := []string{domain}
	input := strings.Join(baseSubdomains, "\n")
	cmd.Stdin = strings.NewReader(input)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("altdns execution failed: %v", err)
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
	RegisterEnumerator(&AltDNSEnumerator{})
}
