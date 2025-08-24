package enumerator

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/itszeeshan/subdomainx/internal/config"
)

type WaybackURLsEnumerator struct{}

func (w *WaybackURLsEnumerator) Name() string {
	return "waybackurls"
}

func (w *WaybackURLsEnumerator) Enumerate(ctx context.Context, domain string, cfg *config.Config, ) ([]string, error) {
	// Build waybackurls command
	// waybackurls expects domains on stdin, so we'll echo the domain and pipe it
	cmd := exec.CommandContext(ctx, "sh", "-c", fmt.Sprintf("echo %s | waybackurls", domain))

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("waybackurls execution failed: %v", err)
	}

	// Parse output (one URL per line)
	lines := strings.Split(string(output), "\n")
	subdomainSet := make(map[string]bool)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Extract subdomain from URL
		// waybackurls returns full URLs, we need to extract the hostname
		if strings.HasPrefix(line, "http://") || strings.HasPrefix(line, "https://") {
			// Remove protocol
			url := strings.TrimPrefix(line, "http://")
			url = strings.TrimPrefix(url, "https://")

			// Get hostname (before first slash or port)
			hostname := url
			if slashIndex := strings.Index(url, "/"); slashIndex != -1 {
				hostname = url[:slashIndex]
			}
			if colonIndex := strings.Index(hostname, ":"); colonIndex != -1 {
				hostname = hostname[:colonIndex]
			}

			// Check if it's a subdomain of our target domain
			if strings.HasSuffix(hostname, "."+domain) && hostname != domain {
				subdomainSet[hostname] = true
			}
		}
	}

	// Convert set to slice
	var subdomains []string
	for subdomain := range subdomainSet {
		subdomains = append(subdomains, subdomain)
	}

	return subdomains, nil
}

func init() {
	RegisterEnumerator(&WaybackURLsEnumerator{})
}
