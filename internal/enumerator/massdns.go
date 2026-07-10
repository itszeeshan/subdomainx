package enumerator

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/itszeeshan/subdomainx/v2/internal/config"
	"github.com/itszeeshan/subdomainx/v2/internal/utils"
)

type MassDNSEnumerator struct{}

func (m *MassDNSEnumerator) Name() string {
	return "massdns"
}

func (m *MassDNSEnumerator) Enumerate(ctx context.Context, domain string, cfg *config.Config) ([]string, error) {
	// Find resolver file
	resolverFile := findResolverFile()
	if resolverFile == "" {
		return nil, fmt.Errorf("massdns: no resolver file found; provide one or install massdns with resolvers")
	}

	// Build subdomain list to resolve
	var subdomains []string
	if cfg.Wordlist != "" {
		words, err := utils.ReadLines(cfg.Wordlist)
		if err != nil {
			return nil, fmt.Errorf("failed to read wordlist: %v", err)
		}
		for _, word := range words {
			subdomains = append(subdomains, fmt.Sprintf("%s.%s", word, domain))
		}
	} else {
		commonWords := []string{"www", "mail", "ftp", "admin", "blog", "dev", "test", "api", "cdn", "static",
			"staging", "app", "ns1", "ns2", "mx", "pop", "imap", "smtp", "vpn", "remote"}
		for _, word := range commonWords {
			subdomains = append(subdomains, fmt.Sprintf("%s.%s", word, domain))
		}
	}

	// Write subdomains to temp file (massdns reads from file more reliably than stdin)
	tmpFile, err := os.CreateTemp("", "massdns-input-*.txt")
	if err != nil {
		return nil, fmt.Errorf("massdns: failed to create temp file: %v", err)
	}
	tmpPath := tmpFile.Name()
	_, _ = tmpFile.WriteString(strings.Join(subdomains, "\n"))
	_ = tmpFile.Close()
	defer func() { _ = os.Remove(tmpPath) }()

	args := []string{"-r", resolverFile, "-t", "A", "-o", "S", tmpPath}
	cmd := exec.CommandContext(ctx, "massdns", args...)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("massdns execution failed: %v", err)
	}

	// Parse massdns output format: subdomain. A IP
	lines := strings.Split(string(output), "\n")
	var results []string
	seen := make(map[string]bool)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 3 && parts[1] == "A" {
			subdomain := strings.TrimSuffix(parts[0], ".")
			if !seen[subdomain] {
				seen[subdomain] = true
				results = append(results, subdomain)
			}
		}
	}

	return results, nil
}

// findResolverFile searches common locations for a DNS resolver list.
func findResolverFile() string {
	candidates := []string{
		"/usr/share/massdns/lists/resolvers.txt",
		"/usr/local/share/massdns/lists/resolvers.txt",
		"/etc/massdns/resolvers.txt",
	}

	if runtime.GOOS == "darwin" {
		// Homebrew paths
		candidates = append([]string{
			"/opt/homebrew/share/massdns/lists/resolvers.txt",
		}, candidates...)
		// Also search Homebrew Cellar
		matches, _ := filepath.Glob("/opt/homebrew/Cellar/massdns/*/.[bB]ottle/etc/lists/resolvers.txt")
		candidates = append(matches, candidates...)
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}

func init() {
	RegisterEnumerator(&MassDNSEnumerator{})
}
