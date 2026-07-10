package tui

import (
	"fmt"
	"runtime"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/itszeeshan/subdomainx/internal/config"
	"github.com/itszeeshan/subdomainx/internal/utils"
)

// RunDashboard launches the Bubble Tea TUI and runs the scan pipeline
// in a background goroutine, routing events through TUIEventSink.
func RunDashboard(cfg *config.Config, checkpoint *utils.Checkpoint, resume string, args []string, outputDir string) error {
	m := initialModel()

	p := tea.NewProgram(m, tea.WithAltScreen())

	sink := NewTUIEventSink(p)

	// Run scan pipeline in background goroutine
	go func() {
		// Send periodic resource stats
		go func() {
			ticker := time.NewTicker(2 * time.Second)
			defer ticker.Stop()
			start := time.Now()
			for range ticker.C {
				var memStats runtime.MemStats
				runtime.ReadMemStats(&memStats)
				p.Send(ResourceMsg{
					MemoryMB:   float64(memStats.Alloc) / 1024 / 1024,
					Goroutines: runtime.NumGoroutine(),
					CPUs:       runtime.NumCPU(),
					Elapsed:    time.Since(start),
				})
			}
		}()

		// Import and run the scan pipeline
		// We call the pipeline executor that main.go sets up
		if scanFunc != nil {
			err := scanFunc(cfg, checkpoint, resume, args, outputDir, sink)
			if err != nil {
				sink.ScanComplete(err)
			}
		} else {
			sink.Log("error", "No scan function registered")
			sink.ScanComplete(fmt.Errorf("no scan function registered"))
		}
	}()

	_, err := p.Run()
	return err
}

// ScanFunc is the type for the scan pipeline function.
type ScanFunc func(cfg *config.Config, checkpoint *utils.Checkpoint, resume string, args []string, outputDir string, sink EventSink) error

// scanFunc is set by the main package to avoid circular imports.
var scanFunc ScanFunc

// RegisterScanFunc allows the main package to register the scan pipeline.
func RegisterScanFunc(f ScanFunc) {
	scanFunc = f
}
