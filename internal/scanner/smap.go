package scanner

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/itszeeshan/subdomainx/internal/config"
	"github.com/itszeeshan/subdomainx/internal/types"
)

type SmapScanner struct{}

func (s *SmapScanner) Name() string {
	return "smap"
}

func (s *SmapScanner) Scan(ctx context.Context, targets []string, cfg *config.Config) ([]types.HTTPResult, error) {
	if len(targets) == 0 {
		return []types.HTTPResult{}, nil
	}

	// Build smap command
	args := []string{
		"-i", "/dev/stdin",
		"-o", "/dev/stdout",
		"-json",
	}

	cmd := exec.CommandContext(ctx, "smap", args...)

	// Create input with targets
	input := strings.Join(targets, "\n")
	cmd.Stdin = strings.NewReader(input)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("smap execution failed: %v", err)
	}

	// Parse JSON output
	lines := strings.Split(string(output), "\n")
	var results []types.HTTPResult

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Parse JSON result
		var smapResult struct {
			Host  string `json:"host"`
			IP    string `json:"ip"`
			Ports []struct {
				Number   int    `json:"number"`
				Protocol string `json:"protocol"`
				State    string `json:"state"`
				Service  string `json:"service"`
				Version  string `json:"version"`
			} `json:"ports"`
		}

		if err := json.Unmarshal([]byte(line), &smapResult); err != nil {
			continue
		}

		// Convert to PortResult
		var ports []types.Port
		for _, port := range smapResult.Ports {
			ports = append(ports, types.Port{
				Number:   port.Number,
				Protocol: port.Protocol,
				State:    port.State,
				Service:  port.Service,
				Version:  port.Version,
			})
		}

		// Convert to HTTPResult for compatibility
		// This is a workaround since the scanner interface expects HTTPResult
		// In a real implementation, you might want to separate these
		httpResult := types.HTTPResult{
			URL:        fmt.Sprintf("http://%s", smapResult.Host),
			StatusCode: 0, // Not applicable for port scan
		}

		results = append(results, httpResult)
	}

	return results, nil
}

func init() {
	RegisterScanner(&SmapScanner{})
}
