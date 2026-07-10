package diff

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/itszeeshan/subdomainx/v2/internal/config"
	"github.com/itszeeshan/subdomainx/v2/internal/types"
)

// Compare computes the diff between the current scan results and a baseline.
// If cfg.BaselineFile is set, it loads from that file; otherwise it finds the
// most recent previous scan from history.
func Compare(cfg *config.Config, scanID string, results []types.SubdomainResult) (*DiffResult, error) {
	current := buildSubdomainMap(results)

	var baseline *HistoryEntry
	var err error

	if cfg.BaselineFile != "" {
		baseline, err = LoadBaselineFromFile(cfg.BaselineFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load baseline: %w", err)
		}
	} else {
		history, err := LoadHistory(cfg.OutputDir)
		if err != nil {
			return nil, fmt.Errorf("no scan history found: %w", err)
		}
		// Find domain from results or config.
		domain := extractDomain(cfg, results)
		baseline = FindBaseline(history, domain, scanID)
		if baseline == nil {
			// First scan — everything is new.
			return firstScanDiff(scanID, current), nil
		}
	}

	return computeDiff(baseline, scanID, current), nil
}

// WriteDiffReport writes the diff result as JSON to the output directory.
func WriteDiffReport(cfg *config.Config, dr *DiffResult) error {
	data, err := json.MarshalIndent(dr, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal diff report: %w", err)
	}

	path := filepath.Join(cfg.OutputDir, cfg.UniqueName+"_diff.json")
	return os.WriteFile(path, data, 0644)
}

// PrintSummary prints a human-readable diff summary to stdout.
func PrintSummary(dr *DiffResult) {
	fmt.Println()
	fmt.Println("--- Diff Summary ---")
	fmt.Printf("Baseline: %s (%s)\n", dr.BaselineScanID, dr.BaselineTime.Format(time.RFC3339))
	fmt.Printf("Current:  %s (%s)\n", dr.CurrentScanID, dr.CurrentTime.Format(time.RFC3339))
	fmt.Println()

	if len(dr.Added) > 0 {
		fmt.Printf("+ %d new subdomains:\n", len(dr.Added))
		for _, s := range dr.Added {
			fmt.Printf("  + %s\n", s)
		}
	}

	if len(dr.Removed) > 0 {
		fmt.Printf("- %d removed subdomains:\n", len(dr.Removed))
		for _, s := range dr.Removed {
			fmt.Printf("  - %s\n", s)
		}
	}

	if len(dr.IPChanges) > 0 {
		fmt.Printf("~ %d IP changes:\n", len(dr.IPChanges))
		for _, c := range dr.IPChanges {
			fmt.Printf("  ~ %s (%v -> %v)\n", c.Subdomain, c.OldIPs, c.NewIPs)
		}
	}

	if len(dr.Added) == 0 && len(dr.Removed) == 0 && len(dr.IPChanges) == 0 {
		fmt.Println("  No changes detected.")
	}

	fmt.Printf("\nTotal: %d current, %d baseline\n", dr.TotalCurrent, dr.TotalBaseline)
	fmt.Println("--------------------")
}

func computeDiff(baseline *HistoryEntry, currentScanID string, current map[string][]string) *DiffResult {
	dr := &DiffResult{
		BaselineScanID: baseline.ScanID,
		BaselineTime:   baseline.Timestamp,
		CurrentScanID:  currentScanID,
		CurrentTime:    time.Now(),
		TotalCurrent:   len(current),
		TotalBaseline:  len(baseline.Subdomains),
	}

	// Find added and IP changes.
	for sub, ips := range current {
		oldIPs, exists := baseline.Subdomains[sub]
		if !exists {
			dr.Added = append(dr.Added, sub)
			continue
		}
		if !ipsEqual(oldIPs, ips) && (len(oldIPs) > 0 || len(ips) > 0) {
			dr.IPChanges = append(dr.IPChanges, IPChange{
				Subdomain: sub,
				OldIPs:    oldIPs,
				NewIPs:    ips,
			})
		}
	}

	// Find removed.
	for sub := range baseline.Subdomains {
		if _, exists := current[sub]; !exists {
			dr.Removed = append(dr.Removed, sub)
		}
	}

	sort.Strings(dr.Added)
	sort.Strings(dr.Removed)
	sort.Slice(dr.IPChanges, func(i, j int) bool {
		return dr.IPChanges[i].Subdomain < dr.IPChanges[j].Subdomain
	})

	return dr
}

func firstScanDiff(scanID string, current map[string][]string) *DiffResult {
	dr := &DiffResult{
		BaselineScanID: "(none)",
		CurrentScanID:  scanID,
		CurrentTime:    time.Now(),
		TotalCurrent:   len(current),
		TotalBaseline:  0,
	}
	for sub := range current {
		dr.Added = append(dr.Added, sub)
	}
	sort.Strings(dr.Added)
	return dr
}

func ipsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	aCopy := make([]string, len(a))
	bCopy := make([]string, len(b))
	copy(aCopy, a)
	copy(bCopy, b)
	sort.Strings(aCopy)
	sort.Strings(bCopy)
	for i := range aCopy {
		if aCopy[i] != bCopy[i] {
			return false
		}
	}
	return true
}

func extractDomain(cfg *config.Config, results []types.SubdomainResult) string {
	// Use UniqueName as a proxy for domain when it was set from the domain arg.
	if cfg.UniqueName != "" && cfg.UniqueName != "scan" {
		return cfg.UniqueName
	}
	// Fallback: extract from first result.
	if len(results) > 0 {
		return results[0].Subdomain
	}
	return "unknown"
}
