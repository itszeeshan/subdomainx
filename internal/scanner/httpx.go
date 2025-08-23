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

type HTTPXScanner struct{}

func (h *HTTPXScanner) Name() string {
	return "httpx"
}

func (h *HTTPXScanner) Scan(ctx context.Context, targets []string, cfg *config.Config) ([]types.HTTPResult, error) {
	if len(targets) == 0 {
		return []types.HTTPResult{}, nil
	}

	// Build httpx command
	args := []string{
		"-l", "/dev/stdin",
		"-o", "/dev/stdout",
		"-json",
		"-title",
		"-tech-detect",
		"-status-code",
		"-content-length",
	}

	cmd := exec.CommandContext(ctx, "httpx", args...)

	// Create input with targets
	input := strings.Join(targets, "\n")
	cmd.Stdin = strings.NewReader(input)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("httpx execution failed: %v", err)
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
		var httpxResult struct {
			URL           string   `json:"url"`
			StatusCode    int      `json:"status-code"`
			Title         string   `json:"title"`
			ContentLength int      `json:"content-length"`
			Technologies  []string `json:"technologies"`
		}

		if err := json.Unmarshal([]byte(line), &httpxResult); err != nil {
			continue
		}

		result := types.HTTPResult{
			URL:           httpxResult.URL,
			StatusCode:    httpxResult.StatusCode,
			Title:         httpxResult.Title,
			ContentLength: httpxResult.ContentLength,
			Technologies:  httpxResult.Technologies,
		}

		results = append(results, result)
	}

	return results, nil
}

func init() {
	RegisterScanner(&HTTPXScanner{})
}
