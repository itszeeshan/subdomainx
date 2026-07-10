package enumerator

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/itszeeshan/subdomainx/internal/config"
)

type AltDNSEnumerator struct{}

func (a *AltDNSEnumerator) Name() string {
	return "altdns"
}

func (a *AltDNSEnumerator) Enumerate(ctx context.Context, domain string, cfg *config.Config) ([]string, error) {
	// Create temp input file with base subdomains
	inputFile, err := os.CreateTemp("", "altdns-input-*.txt")
	if err != nil {
		return nil, fmt.Errorf("altdns: failed to create input file: %v", err)
	}
	inputPath := inputFile.Name()
	_, _ = inputFile.WriteString(domain + "\n")
	_ = inputFile.Close()
	defer func() { _ = os.Remove(inputPath) }()

	// Create temp output file
	outputFile, err := os.CreateTemp("", "altdns-output-*.txt")
	if err != nil {
		return nil, fmt.Errorf("altdns: failed to create output file: %v", err)
	}
	outputPath := outputFile.Name()
	_ = outputFile.Close()
	defer func() { _ = os.Remove(outputPath) }()

	// Find wordlist
	wordlist := cfg.Wordlist
	if wordlist == "" {
		wordlist = findAltdnsWordlist()
		if wordlist == "" {
			return nil, fmt.Errorf("altdns: no wordlist found; use --wordlist to specify one")
		}
	}

	args := []string{"-i", inputPath, "-o", outputPath, "-w", wordlist}
	cmd := exec.CommandContext(ctx, "altdns", args...)
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("altdns execution failed: %v", err)
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("altdns: failed to read output: %v", err)
	}

	lines := strings.Split(string(data), "\n")
	var subdomains []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			subdomains = append(subdomains, line)
		}
	}

	return subdomains, nil
}

// findAltdnsWordlist searches common locations for the altdns wordlist.
func findAltdnsWordlist() string {
	candidates := []string{
		"/usr/share/altdns/words.txt",
		"/usr/local/share/altdns/words.txt",
		"/opt/homebrew/share/altdns/words.txt",
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}

func init() {
	RegisterEnumerator(&AltDNSEnumerator{})
}
