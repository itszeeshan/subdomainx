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

func (s *SmapScanner) Scan(ctx context.Context, targets []string, cfg *config.Config) ([]types.PortResult, error) {
	if len(targets) == 0 {
		return []types.PortResult{}, nil
	}

	// Build smap command
	args := []string{
		"-iL", "-",
		"-oJ", "-",
	}

	cmd := exec.CommandContext(ctx, "smap", args...)

	// Create input with targets
	input := strings.Join(targets, "\n")
	cmd.Stdin = strings.NewReader(input)

	output, err := cmd.Output()
	if err != nil {
		// Try to get stderr for better error reporting
		if exitErr, ok := err.(*exec.ExitError); ok {
			stderr := string(exitErr.Stderr)
			return nil, fmt.Errorf("smap execution failed (exit %d): %s", exitErr.ExitCode(), stderr)
		}
		return nil, fmt.Errorf("smap execution failed: %v", err)
	}

	// Parse JSON output
	lines := strings.Split(string(output), "\n")
	var results []types.PortResult

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

		portResult := types.PortResult{
			Host:  smapResult.Host,
			IP:    smapResult.IP,
			Ports: ports,
		}

		results = append(results, portResult)
	}

	return results, nil
}

func init() {
	RegisterPortScanner(&SmapScanner{})
}
