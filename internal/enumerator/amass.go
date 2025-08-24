package enumerator

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/itszeeshan/subdomainx/internal/config"
)

type AmassEnumerator struct{}

func (a *AmassEnumerator) Name() string {
	return "amass"
}

func (a *AmassEnumerator) Enumerate(ctx context.Context, domain string, cfg *config.Config, ) ([]string, error) {
	// Build amass command with more robust options
	args := []string{
		"enum",
		"-d", domain,
		"-o", "/dev/stdout",
		"-timeout", "60", // Use a reasonable timeout for amass
		"-max-dns-queries", "1000", // Limit DNS queries to avoid rate limiting
		"-passive", // Use passive enumeration first
	}

	// Add wordlist if specified
	if cfg.Wordlist != "" {
		args = append(args, "-w", cfg.Wordlist)
	}

	// Create command with context timeout
	cmd := exec.CommandContext(ctx, "amass", args...)

	// Capture both stdout and stderr
	output, err := cmd.Output()
	if err != nil {
		// Try to get stderr for better error reporting
		if exitErr, ok := err.(*exec.ExitError); ok {
			stderr := string(exitErr.Stderr)
			return nil, fmt.Errorf("amass execution failed (exit %d): %s", exitErr.ExitCode(), stderr)
		}
		return nil, fmt.Errorf("amass execution failed: %v", err)
	}

	// Parse JSON output
	lines := strings.Split(string(output), "\n")
	var subdomains []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var result struct {
			Name string `json:"name"`
		}

		if err := json.Unmarshal([]byte(line), &result); err != nil {
			// Skip malformed JSON lines
			continue
		}

		if result.Name != "" && strings.HasSuffix(result.Name, domain) {
			subdomains = append(subdomains, result.Name)
			// DNS resolution will be done later in the main enumerator
		}
	}

	return subdomains, nil
}

func init() {
	RegisterEnumerator(&AmassEnumerator{})
}
