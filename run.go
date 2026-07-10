package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/itszeeshan/subdomainx/v2/internal/config"
	"github.com/itszeeshan/subdomainx/v2/internal/diff"
	"github.com/itszeeshan/subdomainx/v2/internal/enumerator"
	"github.com/itszeeshan/subdomainx/v2/internal/notify"
	"github.com/itszeeshan/subdomainx/v2/internal/output"
	"github.com/itszeeshan/subdomainx/v2/internal/scanner"
	"github.com/itszeeshan/subdomainx/v2/internal/screenshot"
	"github.com/itszeeshan/subdomainx/v2/internal/server"
	"github.com/itszeeshan/subdomainx/v2/internal/tui"
	"github.com/itszeeshan/subdomainx/v2/internal/types"
	"github.com/itszeeshan/subdomainx/v2/internal/utils"
)

// scanState holds in-progress and completed scan results alongside the
// checkpoint that tracks persistence across interruptions.
type scanState struct {
	checkpoint      *utils.Checkpoint
	results         []types.SubdomainResult
	httpResults     []types.HTTPResult
	portResults     []types.PortResult
	waybackResults  []types.WaybackEntry
	takeoverResults []types.TakeoverResult
}

// initScanState either loads a previous checkpoint (resume mode) or creates a
// fresh one, returning the initial scanState to continue from.
func initScanState(cfg *config.Config, args []string, resume, outputDir string, sink tui.EventSink) (*scanState, error) {
	state := &scanState{}

	if resume != "" {
		cp, err := utils.LoadCheckpoint(resume, outputDir)
		if err != nil {
			return nil, fmt.Errorf("failed to load checkpoint: %v", err)
		}

		sink.Log("info", fmt.Sprintf("Resuming scan from checkpoint: %s", resume))
		sink.Log("info", fmt.Sprintf("Previous progress: %d/%d tasks completed",
			cp.Progress.CompletedTasks, cp.Progress.TotalTasks))
		sink.Log("info", fmt.Sprintf("Previous subdomains found: %d", len(cp.Subdomains)))
		sink.Log("info", fmt.Sprintf("Previous HTTP results: %d", len(cp.HTTPResults)))
		sink.Log("info", fmt.Sprintf("Previous port results: %d", len(cp.PortResults)))

		state.checkpoint = cp
		state.results = cp.Subdomains
		state.httpResults = cp.HTTPResults
		state.portResults = cp.PortResults

		if cp.Domain != "" {
			cfg.WildcardFile = cp.WildcardFile
		}
		return state, nil
	}

	// Fresh scan — create a new checkpoint.
	scanID := cfg.UniqueName
	domain := ""
	if len(args) > 0 {
		domain = args[0]
		if scanID == "scan" {
			scanID = domain
		}
	}

	configMap := map[string]interface{}{
		"threads":   cfg.Threads,
		"retries":   cfg.Retries,
		"timeout":   cfg.Timeout,
		"rateLimit": cfg.RateLimit,
		"wordlist":  cfg.Wordlist,
	}
	state.checkpoint = utils.CreateCheckpoint(scanID, domain, cfg.WildcardFile, configMap)
	return state, nil
}

// executeScanPipeline runs enumeration, optional HTTP scanning, optional port
// scanning, and output generation, persisting progress to the checkpoint after
// each phase.
func executeScanPipeline(cfg *config.Config, state *scanState, resume string, sink tui.EventSink) error {
	cp := state.checkpoint

	// Start signal handler only in CLI mode (TUI and API handle signals differently).
	if _, isTUI := sink.(*tui.TUIEventSink); !isTUI {
		if _, isAPI := sink.(*server.APIEventSink); !isAPI {
			signalHandler := utils.NewSignalHandler(cp, cfg.OutputDir)
			signalHandler.Start()
		}
	}

	// --- Enumeration ---
	if resume == "" || len(state.results) == 0 {
		sink.StageStarted("enumeration", "Starting subdomain enumeration...")
		results, err := enumerator.Run(cfg, sink)
		if err != nil {
			cp.MarkError(fmt.Sprintf("Enumeration failed: %v", err))
			saveCheckpoint(cp, cfg.OutputDir, sink)
			sink.StageCompleted("enumeration", fmt.Sprintf("Enumeration failed: %v", err))
			return fmt.Errorf("enumeration failed: %v", err)
		}
		state.results = results
		cp.AddSubdomains(results)
		cp.UpdateProgress(len(results), len(results))
		saveCheckpoint(cp, cfg.OutputDir, sink)
		sink.StageCompleted("enumeration", fmt.Sprintf("Enumeration completed: %d subdomains", len(results)))
	}

	// --- HTTP scanning ---
	if cfg.Tools["httpx"] && (resume == "" || len(state.httpResults) == 0) {
		sink.StageStarted("http", "Running HTTP scanning with httpx...")
		httpResults, err := scanner.RunHTTPx(cfg, state.results, sink)
		if err != nil {
			sink.Log("error", fmt.Sprintf("HTTP scanning failed: %v", err))
		} else {
			state.httpResults = httpResults
			cp.AddHTTPResults(httpResults)
			saveCheckpoint(cp, cfg.OutputDir, sink)
			sink.HTTPResults(httpResults, len(httpResults))
		}
		sink.StageCompleted("http", fmt.Sprintf("HTTP scanning completed: %d results", len(state.httpResults)))
	}

	// --- Screenshots ---
	if cfg.Screenshot && len(state.httpResults) > 0 {
		sink.StageStarted("screenshot", "Capturing screenshots...")
		count, err := screenshot.CaptureAll(cfg, state.httpResults)
		if err != nil {
			sink.Log("warn", fmt.Sprintf("Screenshot capture failed: %v", err))
		} else {
			sink.Log("info", fmt.Sprintf("Screenshots captured: %d", count))
		}
		sink.StageCompleted("screenshot", "Screenshots done")
	}

	// --- Wayback URLs for HTTP-alive subdomains ---
	if cfg.Tools["waybackurls"] && len(state.httpResults) > 0 && len(state.waybackResults) == 0 {
		sink.StageStarted("wayback", "Collecting Wayback URLs for HTTP-alive subdomains...")
		waybackResults := scanner.RunWaybackURLs(cfg, state.httpResults)
		if len(waybackResults) > 0 {
			totalURLs := 0
			for _, w := range waybackResults {
				totalURLs += len(w.URLs)
			}
			sink.Log("info", fmt.Sprintf("Wayback URLs collected: %d URLs across %d subdomains", totalURLs, len(waybackResults)))
			state.waybackResults = waybackResults
		} else {
			sink.Log("info", "No Wayback URLs found")
		}
		sink.StageCompleted("wayback", "Wayback collection done")
	}

	// --- Technology filter (post-HTTP-scan) ---
	if cfg.TechFilter != "" && len(state.httpResults) > 0 {
		state.httpResults = filterByTechnology(state.httpResults, cfg.TechFilter)
		sink.Log("info", fmt.Sprintf("Filtered HTTP results by technology: %d remaining", len(state.httpResults)))
	}

	// --- Port scanning ---
	if cfg.Tools["smap"] && (resume == "" || len(state.portResults) == 0) {
		sink.StageStarted("ports", "Running port scanning with smap...")
		portResults, err := scanner.RunSmap(cfg, state.results, sink)
		if err != nil {
			sink.Log("error", fmt.Sprintf("Port scanning failed: %v", err))
		} else {
			state.portResults = portResults
			cp.AddPortResults(portResults)
			saveCheckpoint(cp, cfg.OutputDir, sink)
			sink.PortResults(portResults, len(portResults))
		}
		sink.StageCompleted("ports", fmt.Sprintf("Port scanning completed: %d results", len(state.portResults)))
	}

	// --- Subdomain takeover detection ---
	if cfg.Takeover {
		sink.StageStarted("takeover", "Checking for subdomain takeover vulnerabilities...")
		takeoverResults, err := scanner.RunTakeoverCheck(cfg, state.results, state.httpResults, sink)
		if err != nil {
			sink.Log("error", fmt.Sprintf("Takeover detection failed: %v", err))
		} else {
			state.takeoverResults = takeoverResults
			sink.TakeoverResults(takeoverResults)
			if len(takeoverResults) > 0 {
				scanner.PrintTakeoverSummary(takeoverResults, sink)
			}
		}
		sink.StageCompleted("takeover", fmt.Sprintf("Takeover check completed: %d vulnerabilities", len(state.takeoverResults)))
	}

	// --- Record scan history (always, for future diffs) ---
	domain := cp.Domain
	if domain == "" {
		domain = cp.ScanID
	}
	if err := diff.RecordScan(cfg.OutputDir, cp.ScanID, domain, state.results); err != nil {
		sink.Log("warn", fmt.Sprintf("Failed to record scan history: %v", err))
	}

	// --- Diff comparison (before output so HTML report can include diff) ---
	var diffResult *diff.DiffResult
	if cfg.DiffEnabled {
		dr, err := diff.Compare(cfg, cp.ScanID, state.results)
		if err != nil {
			sink.Log("warn", fmt.Sprintf("Diff comparison failed: %v", err))
		} else {
			diffResult = dr
			if err := diff.WriteDiffReport(cfg, dr); err != nil {
				sink.Log("warn", fmt.Sprintf("Failed to write diff report: %v", err))
			}
			diff.PrintSummary(dr)
		}
	}

	// --- Output ---
	sink.StageStarted("output", "Generating output files...")
	if err := output.Generate(cfg, state.results, state.httpResults, state.portResults, state.waybackResults, state.takeoverResults, diffResult); err != nil {
		return fmt.Errorf("failed to generate output: %v", err)
	}
	sink.StageCompleted("output", fmt.Sprintf("Results saved to %s", cfg.OutputDir))

	// --- Notifications ---
	if len(cfg.NotifyChannels) > 0 {
		summary := notify.ScanSummary{
			ScanID:          cp.ScanID,
			Domain:          domain,
			TotalSubdomains: len(state.results),
			TotalHTTP:       len(state.httpResults),
			TotalPorts:      len(state.portResults),
			Duration:        time.Since(cp.Progress.StartTime),
			Diff:            diffResult,
		}
		if err := notify.Send(cfg.NotifyChannels, summary); err != nil {
			sink.Log("warn", fmt.Sprintf("Notification failed: %v", err))
		}
	}

	cp.MarkCompleted()
	saveCheckpoint(cp, cfg.OutputDir, sink)

	sink.ScanComplete(nil)
	return nil
}

// saveCheckpoint persists cp and logs a warning on failure (non-fatal).
func saveCheckpoint(cp *utils.Checkpoint, outputDir string, sink tui.EventSink) {
	if err := utils.SaveCheckpoint(cp, outputDir); err != nil {
		sink.Log("warn", fmt.Sprintf("Failed to save checkpoint: %v", err))
	}
}

// tuiScanFunc bridges the TUI's ScanFunc signature to the internal pipeline.
// It creates a scanState from the checkpoint and runs the full pipeline.
func tuiScanFunc(cfg *config.Config, checkpoint *utils.Checkpoint, resume string, args []string, outputDir string, sink tui.EventSink) error {
	state := &scanState{
		checkpoint:  checkpoint,
		results:     checkpoint.Subdomains,
		httpResults: checkpoint.HTTPResults,
		portResults: checkpoint.PortResults,
	}
	return executeScanPipeline(cfg, state, resume, sink)
}

// filterByTechnology filters HTTP results to only include those matching
// any of the comma-separated technology names.
func filterByTechnology(httpResults []types.HTTPResult, techFilter string) []types.HTTPResult {
	filters := make(map[string]bool)
	for _, t := range strings.Split(techFilter, ",") {
		t = strings.TrimSpace(t)
		if t != "" {
			filters[strings.ToLower(t)] = true
		}
	}
	if len(filters) == 0 {
		return httpResults
	}

	var filtered []types.HTTPResult
	for _, r := range httpResults {
		// Check basic Technologies field
		for _, tech := range r.Technologies {
			if filters[strings.ToLower(strings.TrimSpace(tech))] {
				filtered = append(filtered, r)
				goto next
			}
		}
		// Check DetectedTech field
		for _, dt := range r.DetectedTech {
			if filters[strings.ToLower(dt.Name)] {
				filtered = append(filtered, r)
				goto next
			}
		}
	next:
	}
	return filtered
}
