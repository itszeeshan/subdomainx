package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/itszeeshan/subdomainx/internal/types"
)

const (
	tabDashboard = iota
	tabResults
	tabLogs
	tabCount
)

// toolStatus tracks a single enumeration tool's progress.
type toolStatus struct {
	name   string
	domain string
	status string // "running", "completed", "failed", "skipped"
	found  int
	err    string
}

// model is the Bubble Tea model for the TUI dashboard.
type model struct {
	// Dimensions
	width  int
	height int

	// Tab state
	activeTab int

	// Spinner
	spinner spinner.Model

	// Dashboard state
	tools        []toolStatus
	currentStage string
	stageMessage string
	startTime    time.Time

	// Stats
	totalSubdomains int
	totalHTTP       int
	totalPorts      int
	totalTakeover   int

	// Resource stats
	memoryMB   float64
	goroutines int

	// Results
	subdomains      []types.SubdomainResult
	httpResults     []types.HTTPResult
	portResults     []types.PortResult
	takeoverResults []types.TakeoverResult

	// Results tab state
	resultOffset  int
	sortColumn    int
	sortAsc       bool
	filterText    string
	filterActive  bool

	// Logs tab
	logs        []LogMsg
	logViewport viewport.Model

	// Scan state
	scanDone  bool
	scanError error

	// Quitting
	quitting bool
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	return model{
		activeTab: tabDashboard,
		spinner:   s,
		startTime: time.Now(),
		sortAsc:   true,
		logViewport: viewport.New(80, 20),
	}
}
