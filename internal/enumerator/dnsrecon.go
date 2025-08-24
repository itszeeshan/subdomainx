package enumerator

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/itszeeshan/subdomainx/internal/config"
)

type DNSReconEnumerator struct{}

func (d *DNSReconEnumerator) Name() string {
	return "dnsrecon"
}

func (d *DNSReconEnumerator) Enumerate(ctx context.Context, domain string, cfg *config.Config, ) ([]string, error) {
	// Build dnsrecon command
	args := []string{"-d", domain, "-t", "brt", "--output", "/dev/stdout"}

	cmd := exec.CommandContext(ctx, "dnsrecon", args...)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("dnsrecon execution failed: %v", err)
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
	RegisterEnumerator(&DNSReconEnumerator{})
}
