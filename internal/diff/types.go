package diff

import "time"

// HistoryEntry is the on-disk record for one completed scan.
type HistoryEntry struct {
	ScanID      string              `json:"scan_id"`
	Domain      string              `json:"domain"`
	Timestamp   time.Time           `json:"timestamp"`
	ResultsFile string              `json:"results_file"`
	Subdomains  map[string][]string `json:"subdomains"` // subdomain -> IPs
}

// DiffResult holds the computed differences between two scans.
type DiffResult struct {
	BaselineScanID string     `json:"baseline_scan_id"`
	BaselineTime   time.Time  `json:"baseline_time"`
	CurrentScanID  string     `json:"current_scan_id"`
	CurrentTime    time.Time  `json:"current_time"`
	Added          []string   `json:"added"`
	Removed        []string   `json:"removed"`
	IPChanges      []IPChange `json:"ip_changes,omitempty"`
	TotalCurrent   int        `json:"total_current"`
	TotalBaseline  int        `json:"total_baseline"`
}

// IPChange records a subdomain whose resolved IPs changed between scans.
type IPChange struct {
	Subdomain string   `json:"subdomain"`
	OldIPs    []string `json:"old_ips"`
	NewIPs    []string `json:"new_ips"`
}
