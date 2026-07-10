package server

import (
	"context"
	"sync"
	"time"

	"github.com/itszeeshan/subdomainx/v2/internal/types"
)

// ScanStatus represents the lifecycle of a scan job.
type ScanStatus string

const (
	StatusQueued    ScanStatus = "queued"
	StatusRunning   ScanStatus = "running"
	StatusCompleted ScanStatus = "completed"
	StatusFailed    ScanStatus = "failed"
	StatusCancelled ScanStatus = "cancelled"
)

// ScanJob tracks an individual scan's state, progress, and results.
type ScanJob struct {
	mu          sync.RWMutex
	ID          string         `json:"id"`
	Status      ScanStatus     `json:"status"`
	Domain      string         `json:"domain"`
	Tools       []string       `json:"tools"`
	StartedAt   time.Time      `json:"started_at"`
	CompletedAt *time.Time     `json:"completed_at,omitempty"`
	Progress    ScanProgress   `json:"progress"`
	Results     *ScanResults   `json:"results,omitempty"`
	Error       string         `json:"error,omitempty"`
	cancel      context.CancelFunc
}

// ScanProgress tracks what stage the scan is in.
type ScanProgress struct {
	Stage           string `json:"stage"`
	StageMessage    string `json:"stage_message"`
	SubdomainsFound int    `json:"subdomains_found"`
	HTTPResults     int    `json:"http_results"`
	PortResults     int    `json:"port_results"`
	ToolsCompleted  int    `json:"tools_completed"`
	ToolsFailed     int    `json:"tools_failed"`
}

// ScanResults holds the final output of a completed scan.
type ScanResults struct {
	Subdomains []types.SubdomainResult `json:"subdomains"`
	HTTP       []types.HTTPResult      `json:"http,omitempty"`
	Ports      []types.PortResult      `json:"ports,omitempty"`
	Takeover   []types.TakeoverResult  `json:"takeover,omitempty"`
}

// ScanRequest is the JSON body for POST /api/scan.
type ScanRequest struct {
	Domain     string            `json:"domain"`
	Tools      []string          `json:"tools,omitempty"`
	Threads    int               `json:"threads,omitempty"`
	Retries    int               `json:"retries,omitempty"`
	Timeout    int               `json:"timeout,omitempty"`
	RateLimit  int               `json:"rate_limit,omitempty"`
	Format     string            `json:"format,omitempty"`
	Options    ScanRequestOptions `json:"options,omitempty"`
}

// ScanRequestOptions holds optional feature flags.
type ScanRequestOptions struct {
	Screenshot bool `json:"screenshot,omitempty"`
	TechDetect bool `json:"tech_detect,omitempty"`
	Takeover   bool `json:"takeover,omitempty"`
	Httpx      bool `json:"httpx,omitempty"`
	Smap       bool `json:"smap,omitempty"`
}

// ScanResponse is returned by POST /api/scan (202 Accepted).
type ScanResponse struct {
	ID      string     `json:"id"`
	Status  ScanStatus `json:"status"`
	Message string     `json:"message"`
}

// ErrorResponse is a standard error envelope.
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

// HealthResponse is returned by GET /api/health.
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
	Uptime  string `json:"uptime"`
}

// ScanListItem is a summary for GET /api/scans.
type ScanListItem struct {
	ID          string     `json:"id"`
	Status      ScanStatus `json:"status"`
	Domain      string     `json:"domain"`
	StartedAt   time.Time  `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Subdomains  int        `json:"subdomains"`
}
