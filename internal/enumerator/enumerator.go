package enumerator

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/itszeeshan/subdomainx/internal/cache"
	"github.com/itszeeshan/subdomainx/internal/config"
	"github.com/itszeeshan/subdomainx/internal/types"
	"github.com/itszeeshan/subdomainx/internal/utils"
)

type Enumerator interface {
	Name() string
	Enumerate(ctx context.Context, domain string, cfg *config.Config, cache *cache.DNSCache) ([]string, error)
}

var enumerators = make(map[string]Enumerator)

func RegisterEnumerator(e Enumerator) {
	enumerators[e.Name()] = e
}

func Run(cfg *config.Config, dnsCache *cache.DNSCache) ([]types.SubdomainResult, error) {
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
			fmt.Printf("⚠️  Skipping %s: tool not found in PATH\n", name)
		}
	}

	if len(availableEnumerators) == 0 {
		return nil, fmt.Errorf("no enumeration tools are available")
	}

	fmt.Printf("🔧 Using %d enumeration tools: ", len(availableEnumerators))
	var toolNames []string
	for name := range availableEnumerators {
		toolNames = append(toolNames, name)
	}
	fmt.Printf("%s\n", strings.Join(toolNames, ", "))

	// Use a more reasonable timeout calculation
	totalTimeout := cfg.Timeout * len(domains) * len(availableEnumerators)
	if totalTimeout > 300 { // Cap at 5 minutes
		totalTimeout = 300
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(totalTimeout)*time.Second)
	defer cancel()

	// Create a worker pool for concurrent enumeration
	subdomainChan := make(chan string, 1000)
	var wg sync.WaitGroup

	// Start a goroutine for each enabled and available enumerator and domain
	for _, domain := range domains {
		for name, enumerator := range availableEnumerators {
			wg.Add(1)
			go func(e Enumerator, d string, toolName string) {
				defer wg.Done()

				// Use retry mechanism
				subdomains, err := utils.Retry(func() ([]string, error) {
					return e.Enumerate(ctx, d, cfg, dnsCache)
				}, cfg.Retries, cfg.Timeout)

				if err != nil {
					fmt.Printf("❌ Error with %s on %s: %v\n", toolName, d, err)
					return
				}

				fmt.Printf("✅ %s found %d subdomains for %s\n", toolName, len(subdomains), d)

				// Send subdomains to channel without DNS resolution
				for _, subdomain := range subdomains {
					subdomainChan <- subdomain
				}
			}(enumerator, domain, name)
		}
	}

	// Wait for all enumerators to complete and close subdomain channel
	go func() {
		wg.Wait()
		close(subdomainChan)
	}()

	// Collect and deduplicate subdomains first
	uniqueSubdomains := make(map[string]bool)
	for subdomain := range subdomainChan {
		uniqueSubdomains[subdomain] = true
	}

	fmt.Printf("📊 Total unique subdomains found: %d\n", len(uniqueSubdomains))

	// Now perform DNS resolution only for unique subdomains
	var finalResults []types.SubdomainResult
	for subdomain := range uniqueSubdomains {
		// Resolve DNS for the subdomain
		ips := dnsCache.Resolve(subdomain)
		finalResults = append(finalResults, types.SubdomainResult{
			Subdomain: subdomain,
			Source:    "combined", // Since we deduplicated, we can't track individual sources
			IPs:       ips,
		})
	}

	return finalResults, nil
}
