package enumerator

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"strings"

	"github.com/itszeeshan/subdomainx/v2/internal/config"
)

type FierceEnumerator struct{}

func (f *FierceEnumerator) Name() string {
	return "fierce"
}

func (f *FierceEnumerator) Enumerate(ctx context.Context, domain string, cfg *config.Config) ([]string, error) {
	args := []string{"--domain", domain}

	cmd := exec.CommandContext(ctx, "fierce", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("fierce execution failed: %v", err)
	}

	// Parse fierce output — it prints found subdomains with their IPs
	// Lines look like: "Found: sub.example.com. (1.2.3.4)" or just hostnames
	var subdomains []string
	seen := make(map[string]bool)

	// Match patterns like "Found: hostname.domain. (IP)" or bare hostnames
	foundRe := regexp.MustCompile(`(?i)Found:\s+(\S+)`)
	// Also match lines with IP-to-hostname mappings
	ipHostRe := regexp.MustCompile(`(\S+\.\S+)\.\s+\(?\d`)

	for _, line := range strings.Split(string(output), "\n") {
		line = strings.TrimSpace(line)

		if matches := foundRe.FindStringSubmatch(line); len(matches) > 1 {
			host := strings.TrimSuffix(matches[1], ".")
			if strings.HasSuffix(host, "."+domain) || host == domain {
				if !seen[host] {
					seen[host] = true
					subdomains = append(subdomains, host)
				}
			}
		}

		if matches := ipHostRe.FindStringSubmatch(line); len(matches) > 1 {
			host := strings.TrimSuffix(matches[1], ".")
			if strings.HasSuffix(host, "."+domain) || host == domain {
				if !seen[host] {
					seen[host] = true
					subdomains = append(subdomains, host)
				}
			}
		}

		// Also check for bare domain names in the output
		for _, word := range strings.Fields(line) {
			word = strings.TrimSuffix(word, ".")
			word = strings.Trim(word, "()")
			if strings.HasSuffix(word, "."+domain) && net.ParseIP(word) == nil {
				if !seen[word] {
					seen[word] = true
					subdomains = append(subdomains, word)
				}
			}
		}
	}

	return subdomains, nil
}

func init() {
	RegisterEnumerator(&FierceEnumerator{})
}
