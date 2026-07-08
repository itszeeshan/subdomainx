package main

import (
	"fmt"
	"log"

	"github.com/itszeeshan/subdomainx/internal/config"
	"github.com/itszeeshan/subdomainx/internal/enumerator"
	"github.com/itszeeshan/subdomainx/internal/output"
	"github.com/itszeeshan/subdomainx/internal/scanner"
	"github.com/itszeeshan/subdomainx/internal/types"
	"github.com/itszeeshan/subdomainx/internal/utils"
)

// scanState holds in-progress and completed scan results alongside the
// checkpoint that tracks persistence across interruptions.
type scanState struct {
	checkpoint  *utils.Checkpoint
	results     []types.SubdomainResult
	httpResults []types.HTTPResult
	portResults []types.PortResult
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

	// --- Output ---
	if err := output.Generate(cfg, state.results, state.httpResults, state.portResults); err != nil {
		return fmt.Errorf("failed to generate output: %v", err)
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
