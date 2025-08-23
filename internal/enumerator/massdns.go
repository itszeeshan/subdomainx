package enumerator

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/itszeeshan/subdomainx/internal/cache"
	"github.com/itszeeshan/subdomainx/internal/config"
	"github.com/itszeeshan/subdomainx/internal/utils"
)

type MassDNSEnumerator struct{}

func (m *MassDNSEnumerator) Name() string {
	return "massdns"
}

func (m *MassDNSEnumerator) Enumerate(ctx context.Context, domain string, cfg *config.Config, cache *cache.DNSCache) ([]string, error) {
	// Build massdns command
	args := []string{"-r", "/usr/share/massdns/lists/resolvers.txt", "-t", "A", "-o", "S", "/dev/stdin"}

	cmd := exec.CommandContext(ctx, "massdns", args...)

	// Create input with subdomains
	var subdomains []string
	if cfg.Wordlist != "" {
		// Read wordlist and create subdomains
		words, err := utils.ReadLines(cfg.Wordlist)
		if err != nil {
			return nil, fmt.Errorf("failed to read wordlist: %v", err)
		}

		for _, word := range words {
			subdomain := fmt.Sprintf("%s.%s", word, domain)
			subdomains = append(subdomains, subdomain)
		}
	} else {
		// Use common subdomains
		commonWords := []string{"www", "mail", "ftp", "admin", "blog", "dev", "test", "api", "cdn", "static"}
		for _, word := range commonWords {
			subdomain := fmt.Sprintf("%s.%s", word, domain)
			subdomains = append(subdomains, subdomain)
		}
	}

	// Create input for massdns
	input := strings.Join(subdomains, "\n")
	cmd.Stdin = strings.NewReader(input)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("massdns execution failed: %v", err)
	}

	// Parse massdns output
	lines := strings.Split(string(output), "\n")
	var results []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			// Parse massdns output format: subdomain. A IP
			parts := strings.Fields(line)
			if len(parts) >= 3 && parts[1] == "A" {
				subdomain := parts[0]
				ip := parts[2]

				results = append(results, subdomain)
				// Cache the DNS result with IP
				cache.Store(subdomain, []string{ip})
			}
		}
	}

	return results, nil
}

func init() {
	RegisterEnumerator(&MassDNSEnumerator{})
}
