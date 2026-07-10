package tui

import (
	"time"

	"github.com/itszeeshan/subdomainx/v2/internal/types"
)

// StageMsg signals a pipeline stage transition.
type StageMsg struct {
	Stage   string // "enumeration", "http", "screenshot", "wayback", "ports", "takeover", "output"
	Status  string // "started", "completed", "failed"
	Message string
}

// ToolProgressMsg reports per-tool enumeration progress.
type ToolProgressMsg struct {
	Tool   string
	Domain string
	Status string // "running", "completed", "failed"
	Found  int
	Error  string
}

// ResultMsg delivers discovered subdomains.
type ResultMsg struct {
	Results []types.SubdomainResult
	Total   int
}

// HTTPResultMsg delivers HTTP scan results.
type HTTPResultMsg struct {
	Results []types.HTTPResult
	Total   int
}

// PortResultMsg delivers port scan results.
type PortResultMsg struct {
	Results []types.PortResult
	Total   int
}

// TakeoverResultMsg delivers takeover detection results.
type TakeoverResultMsg struct {
	Results []types.TakeoverResult
}

// LogMsg is a log line to display in the Logs tab.
type LogMsg struct {
	Level   string // "info", "warn", "error"
	Message string
	Time    time.Time
}

// ResourceMsg carries periodic resource usage stats.
type ResourceMsg struct {
	MemoryMB   float64
	Goroutines int
	CPUs       int
	Elapsed    time.Duration
}

// TickMsg triggers periodic resource polling.
type TickMsg time.Time

// ScanCompleteMsg signals the entire scan is done.
type ScanCompleteMsg struct {
	Error error
}
