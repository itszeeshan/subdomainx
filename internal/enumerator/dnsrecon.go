package enumerator

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/itszeeshan/subdomainx/v2/internal/config"
)

type DNSReconEnumerator struct{}

func (d *DNSReconEnumerator) Name() string {
	return "dnsrecon"
}

func (d *DNSReconEnumerator) Enumerate(ctx context.Context, domain string, cfg *config.Config) ([]string, error) {
	// Use JSON output for reliable parsing
	tmpFile, err := os.CreateTemp("", "dnsrecon-*.json")
	if err != nil {
		return nil, fmt.Errorf("dnsrecon: failed to create temp file: %v", err)
	}
	tmpPath := tmpFile.Name()
	_ = tmpFile.Close()
	defer func() { _ = os.Remove(tmpPath) }()

	// Use standard enumeration (-t std) which does SOA, NS, MX, A, AAAA, TXT, zone transfer
	args := []string{"-d", domain, "-t", "std", "-j", tmpPath}

	cmd := exec.CommandContext(ctx, "dnsrecon", args...)
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("dnsrecon execution failed: %v", err)
	}

	// Parse JSON output
	data, err := os.ReadFile(tmpPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("dnsrecon: failed to read output: %v", err)
	}

	var records []map[string]interface{}
	if err := json.Unmarshal(data, &records); err != nil {
		// Try to extract subdomains from raw output instead
		return extractSubdomainsFromText(string(data), domain), nil
	}

	var subdomains []string
	seen := make(map[string]bool)

	for _, record := range records {
		// Check "name" field which contains the hostname
		if name, ok := record["name"].(string); ok {
			name = strings.TrimSuffix(name, ".")
			if (strings.HasSuffix(name, "."+domain) || name == domain) && !seen[name] {
				seen[name] = true
				subdomains = append(subdomains, name)
			}
		}
		// Also check "exchange" field (MX records)
		if exchange, ok := record["exchange"].(string); ok {
			exchange = strings.TrimSuffix(exchange, ".")
			if strings.HasSuffix(exchange, "."+domain) && !seen[exchange] {
				seen[exchange] = true
				subdomains = append(subdomains, exchange)
			}
		}
		// Check "target" field (CNAME, NS records)
		if target, ok := record["target"].(string); ok {
			target = strings.TrimSuffix(target, ".")
			if strings.HasSuffix(target, "."+domain) && !seen[target] {
				seen[target] = true
				subdomains = append(subdomains, target)
			}
		}
	}

	return subdomains, nil
}

// extractSubdomainsFromText is a fallback parser that finds domain names in text output.
func extractSubdomainsFromText(text, domain string) []string {
	var subdomains []string
	seen := make(map[string]bool)

	for _, line := range strings.Split(text, "\n") {
		for _, word := range strings.Fields(line) {
			word = strings.TrimSuffix(word, ".")
			word = strings.Trim(word, "[](),\"")
			if strings.HasSuffix(word, "."+domain) && !seen[word] {
				seen[word] = true
				subdomains = append(subdomains, word)
			}
		}
	}
	return subdomains
}

func init() {
	RegisterEnumerator(&DNSReconEnumerator{})
}
