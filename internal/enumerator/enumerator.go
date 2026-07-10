package enumerator

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/itszeeshan/subdomainx/internal/config"
	"github.com/itszeeshan/subdomainx/internal/tui"
	"github.com/itszeeshan/subdomainx/internal/types"
	"github.com/itszeeshan/subdomainx/internal/utils"
)

type Enumerator interface {
	Name() string
	Enumerate(ctx context.Context, domain string, cfg *config.Config) ([]string, error)
}

var enumerators = make(map[string]Enumerator)

func RegisterEnumerator(e Enumerator) {
	enumerators[e.Name()] = e
}

func Run(cfg *config.Config, sink tui.EventSink) ([]types.SubdomainResult, error) {
	// Read domains from wildcard file
	domains, err := utils.ReadLines(cfg.WildcardFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read wildcard file: %v", err)
	}

	// Filter enumerators to only include available tools
	availableEnumerators := make(map[string]Enumerator)
	for name, enumerator := range enumerators {
		if cfg.Tools[name] && utils.CheckToolAvailability(name) {
			availableEnumerators[name] = enumerator
		} else if cfg.Tools[name] && !utils.CheckToolAvailability(name) {
			sink.ToolProgress(name, "", "skipped", 0, nil)
		}
	}

	if len(availableEnumerators) == 0 {
		return nil, fmt.Errorf("no enumeration tools are available")
	}

	var toolNames []string
	for name := range availableEnumerators {
		toolNames = append(toolNames, name)
	}
	sink.Log("info", fmt.Sprintf("Using %d enumeration tools: %s", len(availableEnumerators), strings.Join(toolNames, ", ")))

	// Use a more reasonable timeout calculation
	totalTimeout := cfg.Timeout * len(domains) * len(availableEnumerators)
	if totalTimeout > 300 { // Cap at 5 minutes
		totalTimeout = 300
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(totalTimeout)*time.Second)
	defer cancel()

	type subdomainEntry struct {
		subdomain string
		source    string
	}

	// Create a worker pool for concurrent enumeration
	subdomainChan := make(chan subdomainEntry, 1000)
	var wg sync.WaitGroup

	// Calculate total tasks for progress tracking
	totalTasks := len(domains) * len(availableEnumerators)
	completedTasks := 0
	var progressMutex sync.Mutex

	// Start progress tracking
	utils.StartEnumerationProgress(totalTasks)

	// Start a goroutine for each enabled and available enumerator and domain
	for _, domain := range domains {
		for name, enumerator := range availableEnumerators {
			wg.Add(1)
			go func(e Enumerator, d string, toolName string) {
				defer wg.Done()

				// Use retry mechanism
				subdomains, err := utils.Retry(func() ([]string, error) {
					return e.Enumerate(ctx, d, cfg)
				}, cfg.Retries, cfg.Timeout)

				if err != nil {
					sink.ToolProgress(toolName, d, "failed", 0, err)
				} else {
					sink.ToolProgress(toolName, d, "completed", len(subdomains), nil)
				}

				// Update progress
				progressMutex.Lock()
				completedTasks++
				utils.UpdateEnumerationProgress(completedTasks)
				progressMutex.Unlock()

				// Send subdomains to channel without DNS resolution
				for _, subdomain := range subdomains {
					subdomainChan <- subdomainEntry{subdomain: subdomain, source: toolName}
				}
			}(enumerator, domain, name)
		}
	}

	// Wait for all enumerators to complete and close subdomain channel
	go func() {
		wg.Wait()
		utils.FinishEnumerationProgress()
		close(subdomainChan)
	}()

	// Periodic resource check during enumeration
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				utils.CheckResources()
			case <-ctx.Done():
				return
			}
		}
	}()

	// Collect and deduplicate subdomains, tracking sources per subdomain
	subdomainSources := make(map[string][]string)
	for entry := range subdomainChan {
		sources := subdomainSources[entry.subdomain]
		alreadyHas := false
		for _, s := range sources {
			if s == entry.source {
				alreadyHas = true
				break
			}
		}
		if !alreadyHas {
			subdomainSources[entry.subdomain] = append(subdomainSources[entry.subdomain], entry.source)
		}
	}

	// Create results without DNS resolution for better performance
	var finalResults []types.SubdomainResult
	for subdomain, sources := range subdomainSources {
		finalResults = append(finalResults, types.SubdomainResult{
			Subdomain: subdomain,
			Source:    strings.Join(sources, ","),
			IPs:       []string{},
		})
	}

	sink.SubdomainsFound(finalResults, len(subdomainSources))

	return finalResults, nil
}
