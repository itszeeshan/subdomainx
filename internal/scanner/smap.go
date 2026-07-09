package scanner

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/itszeeshan/subdomainx/internal/config"
	"github.com/itszeeshan/subdomainx/internal/types"
)

type SmapScanner struct{}

func (s *SmapScanner) Name() string {
	return "smap"
}

// smapEntry matches smap's actual JSON output format.
type smapEntry struct {
	IP           string `json:"ip"`
	UserHostname string `json:"user_hostname"`
	Ports        []struct {
		Port     int    `json:"port"`
		Service  string `json:"service"`
		Protocol string `json:"protocol"`
	} `json:"ports"`
	OS struct {
		Name string `json:"name"`
	} `json:"os"`
}

func (s *SmapScanner) Scan(ctx context.Context, targets []string, cfg *config.Config) ([]types.PortResult, error) {
	if len(targets) == 0 {
		return []types.PortResult{}, nil
	}

	// smap requires a real file for -oJ, it can't write JSON to stdout.
	tmpDir := os.TempDir()
	outFile := filepath.Join(tmpDir, fmt.Sprintf("subdomainx-smap-%d.json", os.Getpid()))
	defer func() { _ = os.Remove(outFile) }()

	args := []string{"-iL", "-", "-oJ", outFile}
	cmd := exec.CommandContext(ctx, "smap", args...)
	cmd.Stdin = strings.NewReader(strings.Join(targets, "\n"))

	output, err := cmd.CombinedOutput()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("smap execution failed (exit %d): %s", exitErr.ExitCode(), string(output))
		}
		return nil, fmt.Errorf("smap execution failed: %v", err)
	}

	// Read JSON output file
	data, err := os.ReadFile(outFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read smap output: %v", err)
	}

	// smap outputs a JSON array
	var entries []smapEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("failed to parse smap JSON: %v", err)
	}

	var results []types.PortResult
	for _, entry := range entries {
		var ports []types.Port
		for _, p := range entry.Ports {
			service := strings.TrimSuffix(p.Service, "?")
			ports = append(ports, types.Port{
				Number:   p.Port,
				Protocol: p.Protocol,
				State:    "open",
				Service:  service,
			})
		}
		if len(ports) == 0 {
			continue
		}
		host := entry.UserHostname
		if host == "" {
			host = entry.IP
		}
		results = append(results, types.PortResult{
			Host:  host,
			IP:    entry.IP,
			Ports: ports,
		})
	}

	return results, nil
}

func init() {
	RegisterPortScanner(&SmapScanner{})
}
