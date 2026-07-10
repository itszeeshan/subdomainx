package diff

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/itszeeshan/subdomainx/v2/internal/types"
)

const (
	historyFileName   = ".scan_history.json"
	maxEntriesPerDomain = 100
)

// RecordScan appends a new entry to the scan history file. It prunes old
// entries to keep at most maxEntriesPerDomain per domain.
func RecordScan(outputDir, scanID, domain string, results []types.SubdomainResult) error {
	history, _ := LoadHistory(outputDir) // ignore error — start fresh if missing

	entry := HistoryEntry{
		ScanID:    scanID,
		Domain:    domain,
		Timestamp: time.Now(),
		Subdomains: buildSubdomainMap(results),
	}
	history = append(history, entry)
	history = pruneHistory(history, domain)

	return saveHistory(outputDir, history)
}

// LoadHistory reads the scan history from disk.
func LoadHistory(outputDir string) ([]HistoryEntry, error) {
	path := filepath.Join(outputDir, historyFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var history []HistoryEntry
	if err := json.Unmarshal(data, &history); err != nil {
		return nil, fmt.Errorf("failed to parse scan history: %w", err)
	}
	return history, nil
}

// FindBaseline returns the most recent history entry for the given domain,
// excluding the entry with excludeScanID (the current scan).
func FindBaseline(history []HistoryEntry, domain, excludeScanID string) *HistoryEntry {
	var best *HistoryEntry
	for i := range history {
		e := &history[i]
		if e.Domain != domain || e.ScanID == excludeScanID {
			continue
		}
		if best == nil || e.Timestamp.After(best.Timestamp) {
			best = e
		}
	}
	return best
}

// LoadBaselineFromFile loads a baseline from a results JSON file (the format
// produced by output.Generate) or from a history entry file.
func LoadBaselineFromFile(path string) (*HistoryEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read baseline file: %w", err)
	}

	// Try as ScanResults first (output.Generate format).
	var scanResults types.ScanResults
	if err := json.Unmarshal(data, &scanResults); err == nil && len(scanResults.Subdomains) > 0 {
		entry := &HistoryEntry{
			ScanID:     "baseline",
			Timestamp:  time.Now(),
			Subdomains: buildSubdomainMap(scanResults.Subdomains),
		}
		return entry, nil
	}

	// Try as a raw array of SubdomainResult.
	var subResults []types.SubdomainResult
	if err := json.Unmarshal(data, &subResults); err == nil && len(subResults) > 0 {
		entry := &HistoryEntry{
			ScanID:     "baseline",
			Timestamp:  time.Now(),
			Subdomains: buildSubdomainMap(subResults),
		}
		return entry, nil
	}

	// Try as HistoryEntry directly.
	var entry HistoryEntry
	if err := json.Unmarshal(data, &entry); err == nil && len(entry.Subdomains) > 0 {
		return &entry, nil
	}

	return nil, fmt.Errorf("unrecognized baseline file format: %s", path)
}

func buildSubdomainMap(results []types.SubdomainResult) map[string][]string {
	m := make(map[string][]string, len(results))
	for _, r := range results {
		m[r.Subdomain] = r.IPs
	}
	return m
}

func pruneHistory(history []HistoryEntry, domain string) []HistoryEntry {
	// Separate entries for this domain from others.
	var domainEntries []HistoryEntry
	var others []HistoryEntry
	for _, e := range history {
		if e.Domain == domain {
			domainEntries = append(domainEntries, e)
		} else {
			others = append(others, e)
		}
	}

	if len(domainEntries) <= maxEntriesPerDomain {
		return history
	}

	// Sort by timestamp descending, keep only the newest.
	sort.Slice(domainEntries, func(i, j int) bool {
		return domainEntries[i].Timestamp.After(domainEntries[j].Timestamp)
	})
	domainEntries = domainEntries[:maxEntriesPerDomain]

	return append(others, domainEntries...)
}

func saveHistory(outputDir string, history []HistoryEntry) error {
	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal scan history: %w", err)
	}

	path := filepath.Join(outputDir, historyFileName)
	tmpPath := path + ".tmp"

	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write scan history: %w", err)
	}

	return os.Rename(tmpPath, path)
}
