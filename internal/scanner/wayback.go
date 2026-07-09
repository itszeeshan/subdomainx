package scanner

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"

	"github.com/itszeeshan/subdomainx/internal/config"
	"github.com/itszeeshan/subdomainx/internal/types"
	"github.com/itszeeshan/subdomainx/internal/utils"
)

// RunWaybackURLs collects historical URLs for each HTTP-alive subdomain.
func RunWaybackURLs(cfg *config.Config, httpResults []types.HTTPResult) []types.WaybackEntry {
	if len(httpResults) == 0 {
		return nil
	}

	// Check if waybackurls is installed
	if _, err := exec.LookPath("waybackurls"); err != nil {
		log.Println("Warning: waybackurls not installed, skipping wayback collection")
		return nil
	}

	// Deduplicate hostnames from HTTP results
	hostMap := make(map[string]string) // hostname -> parent domain
	for _, h := range httpResults {
		host := h.URL
		for _, prefix := range []string{"https://", "http://"} {
			host = strings.TrimPrefix(host, prefix)
		}
		if idx := strings.Index(host, "/"); idx != -1 {
			host = host[:idx]
		}
		if idx := strings.Index(host, ":"); idx != -1 {
			host = host[:idx]
		}
		if host != "" {
			parts := strings.Split(host, ".")
			domain := host
			if len(parts) > 2 {
				domain = strings.Join(parts[len(parts)-2:], ".")
			}
			hostMap[host] = domain
		}
	}

	pool := utils.NewWorkerPool(cfg.Threads, cfg.RateLimit)
	defer pool.Stop()

	var mu sync.Mutex
	var results []types.WaybackEntry
	var wg sync.WaitGroup

	for host, domain := range hostMap {
		host, domain := host, domain
		wg.Add(1)
		pool.Submit(func() {
			defer wg.Done()
			urls := fetchWaybackURLs(host, cfg.Timeout)
			if len(urls) > 0 {
				mu.Lock()
				results = append(results, types.WaybackEntry{
					Subdomain: host,
					Domain:    domain,
					URLs:      urls,
				})
				mu.Unlock()
			}
		})
	}

	wg.Wait()
	return results
}

func fetchWaybackURLs(host string, timeout int) []string {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("echo %s | waybackurls", host))
	output, err := cmd.Output()
	if err != nil {
		return nil
	}

	seen := make(map[string]bool)
	var urls []string
	for _, line := range strings.Split(string(output), "\n") {
		line = strings.TrimSpace(line)
		if line != "" && !seen[line] {
			seen[line] = true
			urls = append(urls, line)
		}
	}
	return urls
}
