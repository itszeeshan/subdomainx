package enumerator

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/itszeeshan/subdomainx/internal/config"
)

type FindomainEnumerator struct{}

func (f *FindomainEnumerator) Name() string {
	return "findomain"
}

func (f *FindomainEnumerator) Enumerate(ctx context.Context, domain string, cfg *config.Config) ([]string, error) {
	// Use a temp file for output since findomain writes to file, not stdout
	tmpFile, err := os.CreateTemp("", "findomain-*.txt")
	if err != nil {
		return nil, fmt.Errorf("findomain: failed to create temp file: %v", err)
	}
	tmpPath := tmpFile.Name()
	_ = tmpFile.Close()
	defer func() { _ = os.Remove(tmpPath) }()

	args := []string{"-t", domain, "-u", tmpPath}

	cmd := exec.CommandContext(ctx, "findomain", args...)
	cmd.Stderr = nil // suppress stderr noise

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("findomain execution failed: %v", err)
	}

	data, err := os.ReadFile(tmpPath)
	if err != nil {
		// File may not exist if no results found
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("findomain: failed to read output: %v", err)
	}

	lines := strings.Split(string(data), "\n")
	var subdomains []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			subdomains = append(subdomains, line)
		}
	}

	// Clean up the default output file findomain sometimes creates
	defaultOut := filepath.Join(".", domain+".txt")
	_ = os.Remove(defaultOut)

	return subdomains, nil
}

func init() {
	RegisterEnumerator(&FindomainEnumerator{})
}
