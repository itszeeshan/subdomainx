package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/itszeeshan/subdomainx/internal/config"
	"github.com/itszeeshan/subdomainx/internal/diff"
	"github.com/itszeeshan/subdomainx/internal/enumerator"
	"github.com/itszeeshan/subdomainx/internal/notify"
	"github.com/itszeeshan/subdomainx/internal/output"
	"github.com/itszeeshan/subdomainx/internal/screenshot"
	"github.com/itszeeshan/subdomainx/internal/scanner"
	"github.com/itszeeshan/subdomainx/internal/types"
	"github.com/itszeeshan/subdomainx/internal/utils"
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
func initScanState(cfg *config.Config, args []string, resume, outputDir string) (*scanState, error) {
	state := &scanState{}

	if resume != "" {
		cp, err := utils.LoadCheckpoint(resume, outputDir)
		if err != nil {
			return nil, fmt.Errorf("failed to load checkpoint: %v", err)
		}

		fmt.Printf("🔄 Resuming scan from checkpoint: %s\n", resume)
		fmt.Printf("📊 Previous progress: %d/%d tasks completed\n",
			cp.Progress.CompletedTasks, cp.Progress.TotalTasks)
		fmt.Printf("🔍 Previous subdomains found: %d\n", len(cp.Subdomains))
		fmt.Printf("🌐 Previous HTTP results: %d\n", len(cp.HTTPResults))
		fmt.Printf("🔌 Previous port results: %d\n", len(cp.PortResults))

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
func executeScanPipeline(cfg *config.Config, state *scanState, resume string) error {
	cp := state.checkpoint

	// Start signal handler so Ctrl-C saves progress gracefully.
	signalHandler := utils.NewSignalHandler(cp, cfg.OutputDir)
	signalHandler.Start()

	// --- Enumeration ---
	if resume == "" || len(state.results) == 0 {
		results, err := enumerator.Run(cfg)
		if err != nil {
			cp.MarkError(fmt.Sprintf("Enumeration failed: %v", err))
			saveCheckpoint(cp, cfg.OutputDir)
			return fmt.Errorf("enumeration failed: %v", err)
		}
		state.results = results
		cp.AddSubdomains(results)
		cp.UpdateProgress(len(results), len(results))
		saveCheckpoint(cp, cfg.OutputDir)
	}

	// --- HTTP scanning ---
	if cfg.Tools["httpx"] && (resume == "" || len(state.httpResults) == 0) {
		log.Println("🔍 Running HTTP scanning with httpx...")
		httpResults, err := scanner.RunHTTPx(cfg, state.results)
		if err != nil {
			log.Printf("HTTP scanning failed: %v", err)
		} else {
			log.Printf("✅ HTTP scanning completed: %d results", len(httpResults))
			state.httpResults = httpResults
			cp.AddHTTPResults(httpResults)
			saveCheckpoint(cp, cfg.OutputDir)
		}
	}

	// --- Screenshots ---
	if cfg.Screenshot && len(state.httpResults) > 0 {
		log.Println("📸 Capturing screenshots...")
		count, err := screenshot.CaptureAll(cfg, state.httpResults)
		if err != nil {
			log.Printf("Warning: Screenshot capture failed: %v", err)
		} else {
			log.Printf("✅ Screenshots captured: %d", count)
		}
	}

	// --- Wayback URLs for HTTP-alive subdomains ---
	if cfg.Tools["waybackurls"] && len(state.httpResults) > 0 && len(state.waybackResults) == 0 {
		log.Println("🕰️  Collecting Wayback URLs for HTTP-alive subdomains...")
		waybackResults := scanner.RunWaybackURLs(cfg, state.httpResults)
		if len(waybackResults) > 0 {
			totalURLs := 0
			for _, w := range waybackResults {
				totalURLs += len(w.URLs)
			}
			log.Printf("✅ Wayback URLs collected: %d URLs across %d subdomains", totalURLs, len(waybackResults))
			state.waybackResults = waybackResults
		} else {
			log.Println("ℹ️  No Wayback URLs found")
		}
	}

	// --- Technology filter (post-HTTP-scan) ---
	if cfg.TechFilter != "" && len(state.httpResults) > 0 {
		state.httpResults = filterByTechnology(state.httpResults, cfg.TechFilter)
		log.Printf("🔧 Filtered HTTP results by technology: %d remaining", len(state.httpResults))
	}

	// --- Port scanning ---
	if cfg.Tools["smap"] && (resume == "" || len(state.portResults) == 0) {
		log.Println("🔍 Running port scanning with smap...")
		portResults, err := scanner.RunSmap(cfg, state.results)
		if err != nil {
			log.Printf("Port scanning failed: %v", err)
		} else {
			log.Printf("✅ Port scanning completed: %d results", len(portResults))
			state.portResults = portResults
			cp.AddPortResults(portResults)
			saveCheckpoint(cp, cfg.OutputDir)
		}
	}

	// --- Subdomain takeover detection ---
	if cfg.Takeover {
		log.Println("🔍 Checking for subdomain takeover vulnerabilities...")
		takeoverResults, err := scanner.RunTakeoverCheck(cfg, state.results, state.httpResults)
		if err != nil {
			log.Printf("Takeover detection failed: %v", err)
		} else {
			log.Printf("✅ Takeover check completed: %d potential vulnerabilities", len(takeoverResults))
			state.takeoverResults = takeoverResults
			if len(takeoverResults) > 0 {
				scanner.PrintTakeoverSummary(takeoverResults)
			}
		}
	}

	// --- Record scan history (always, for future diffs) ---
	domain := cp.Domain
	if domain == "" {
		domain = cp.ScanID
	}
	if err := diff.RecordScan(cfg.OutputDir, cp.ScanID, domain, state.results); err != nil {
		log.Printf("Warning: Failed to record scan history: %v", err)
	}

	// --- Diff comparison (before output so HTML report can include diff) ---
	var diffResult *diff.DiffResult
	if cfg.DiffEnabled {
		dr, err := diff.Compare(cfg, cp.ScanID, state.results)
		if err != nil {
			log.Printf("Warning: Diff comparison failed: %v", err)
		} else {
			diffResult = dr
			if err := diff.WriteDiffReport(cfg, dr); err != nil {
				log.Printf("Warning: Failed to write diff report: %v", err)
			}
			diff.PrintSummary(dr)
		}
	}

	// --- Output ---
	if err := output.Generate(cfg, state.results, state.httpResults, state.portResults, state.waybackResults, state.takeoverResults, diffResult); err != nil {
		return fmt.Errorf("failed to generate output: %v", err)
	}

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
			log.Printf("Warning: %v", err)
		}
	}

	cp.MarkCompleted()
	saveCheckpoint(cp, cfg.OutputDir)

	log.Printf("Scan completed. Results saved to %s", cfg.OutputDir)
	return nil
}

// saveCheckpoint persists cp and logs a warning on failure (non-fatal).
func saveCheckpoint(cp *utils.Checkpoint, outputDir string) {
	if err := utils.SaveCheckpoint(cp, outputDir); err != nil {
		log.Printf("Warning: Failed to save checkpoint: %v", err)
	}
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
