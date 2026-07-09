package scanner

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/itszeeshan/subdomainx/internal/config"
	"github.com/itszeeshan/subdomainx/internal/types"
	"github.com/itszeeshan/subdomainx/internal/utils"
)

// Scanner interface for different scanning types
type Scanner interface {
	Name() string
	Scan(ctx context.Context, targets []string, cfg *config.Config) ([]types.HTTPResult, error)
}

// PortScanner interface for port scanning
type PortScanner interface {
	Name() string
	Scan(ctx context.Context, targets []string, cfg *config.Config) ([]types.PortResult, error)
}

var scanners = make(map[string]Scanner)
var portScanners = make(map[string]PortScanner)

// RegisterScanner registers a new scanner
func RegisterScanner(s Scanner) {
	scanners[s.Name()] = s
}

// RegisterPortScanner registers a new port scanner
func RegisterPortScanner(s PortScanner) {
	portScanners[s.Name()] = s
}

// RunHTTPx runs HTTP scanning on discovered subdomains
func RunHTTPx(cfg *config.Config, subdomains []types.SubdomainResult) ([]types.HTTPResult, error) {
	if len(subdomains) == 0 {
		return []types.HTTPResult{}, nil
	}

	// Convert subdomains to URLs with limit for performance
	var urls []string
	maxTargets := cfg.MaxHTTPTargets // Use configurable limit

	for i, subdomain := range subdomains {
		if i >= maxTargets {
			fmt.Printf("⚠️  Limiting HTTP scan to first %d subdomains for performance\n", maxTargets)
			break
		}
		urls = append(urls, fmt.Sprintf("http://%s", subdomain.Subdomain))
		urls = append(urls, fmt.Sprintf("https://%s", subdomain.Subdomain))
	}

	// Use httpx scanner if available
	if httpxScanner, exists := scanners["httpx"]; exists {
		results, err := httpxScanner.Scan(context.Background(), urls, cfg)
		if err != nil {
			return nil, err
		}
		// When tech detection is enabled, enrich results with fingerprinting
		if cfg.TechDetect {
			enrichWithFingerprinting(cfg, results)
		}
		return results, nil
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout*len(urls))*time.Second)
	defer cancel()

	// Create worker pool
	pool := utils.NewWorkerPool(cfg.Threads, cfg.RateLimit)
	defer pool.Stop()

	// Create channels for results and errors
	results := make(chan types.HTTPResult, len(urls))
	errors := make(chan error, len(urls))
	var wg sync.WaitGroup

	// Submit scanning jobs
	for _, url := range urls {
		url := url // capture loop variable for goroutine
		wg.Add(1)
		pool.Submit(func() {
			defer wg.Done()

			// Use retry mechanism
			result, err := utils.Retry(func() (types.HTTPResult, error) {
				return scanHTTP(ctx, url, cfg)
			}, cfg.Retries, cfg.Timeout)

			if err != nil {
				errors <- fmt.Errorf("failed to scan %s: %v", url, err)
				return
			}

			// Apply filters
			if shouldIncludeHTTPResult(result, cfg) {
				results <- result
			}
		})
	}

	// Wait for all jobs to complete
	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	// Collect results
	var httpResults []types.HTTPResult
	for result := range results {
		httpResults = append(httpResults, result)
	}

	// Log errors (but don't fail the scan)
	for err := range errors {
		fmt.Printf("HTTP scan error: %v\n", err)
	}

	return httpResults, nil
}

// RunSmap runs port scanning on discovered subdomains
func RunSmap(cfg *config.Config, subdomains []types.SubdomainResult) ([]types.PortResult, error) {
	if len(subdomains) == 0 {
		return []types.PortResult{}, nil
	}

	// Extract unique hosts
	hosts := make(map[string]bool)
	for _, subdomain := range subdomains {
		hosts[subdomain.Subdomain] = true
	}

	var uniqueHosts []string
	for host := range hosts {
		uniqueHosts = append(uniqueHosts, host)
	}

	// Use smap scanner if available
	if smapScanner, exists := portScanners["smap"]; exists {
		return smapScanner.Scan(context.Background(), uniqueHosts, cfg)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout*len(uniqueHosts))*time.Second)
	defer cancel()

	// Create worker pool
	pool := utils.NewWorkerPool(cfg.Threads, cfg.RateLimit)
	defer pool.Stop()

	// Create channels for results and errors
	results := make(chan types.PortResult, len(uniqueHosts))
	errors := make(chan error, len(uniqueHosts))
	var wg sync.WaitGroup

	// Submit scanning jobs
	for _, host := range uniqueHosts {
		host := host // capture loop variable for goroutine
		wg.Add(1)
		pool.Submit(func() {
			defer wg.Done()

			// Use retry mechanism
			result, err := utils.Retry(func() (types.PortResult, error) {
				return scanPorts(ctx, host, cfg)
			}, cfg.Retries, cfg.Timeout)

			if err != nil {
				errors <- fmt.Errorf("failed to scan ports for %s: %v", host, err)
				return
			}

			// Apply filters
			if shouldIncludePortResult(result, cfg) {
				results <- result
			}
		})
	}

	// Wait for all jobs to complete
	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	// Collect results
	var portResults []types.PortResult
	for result := range results {
		portResults = append(portResults, result)
	}

	// Log errors (but don't fail the scan)
	for err := range errors {
		fmt.Printf("Port scan error: %v\n", err)
	}

	return portResults, nil
}

// shouldIncludeHTTPResult checks if an HTTP result should be included based on filters
func shouldIncludeHTTPResult(result types.HTTPResult, cfg *config.Config) bool {
	if statusCodes, exists := cfg.Filters["status_code"]; exists && statusCodes != "" {
		resultCode := strconv.Itoa(result.StatusCode)
		for _, code := range strings.Split(statusCodes, ",") {
			if strings.TrimSpace(code) == resultCode {
				return true
			}
		}
		return false
	}
	return true
}

// shouldIncludePortResult checks if a port result should be included based on filters
func shouldIncludePortResult(result types.PortResult, cfg *config.Config) bool {
	if ports, exists := cfg.Filters["ports"]; exists && ports != "" {
		portSet := make(map[string]bool)
		for _, p := range strings.Split(ports, ",") {
			portSet[strings.TrimSpace(p)] = true
		}
		for _, port := range result.Ports {
			if portSet[strconv.Itoa(port.Number)] {
				return true
			}
		}
		return false
	}
	return true
}

// scanHTTP performs HTTP scanning on a single URL
func scanHTTP(ctx context.Context, url string, cfg *config.Config) (types.HTTPResult, error) {
	client := &http.Client{
		Timeout: time.Duration(cfg.Timeout) * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return types.HTTPResult{}, err
	}

	// Set user agent
	req.Header.Set("User-Agent", "SubdomainX/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return types.HTTPResult{}, err
	}
	defer func() { _ = resp.Body.Close() }()

	// Read body once for title extraction and fingerprinting
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 64*1024))

	title := extractTitleFromBody(body)
	technologies := extractTechnologies(resp)

	result := types.HTTPResult{
		URL:           url,
		StatusCode:    resp.StatusCode,
		Title:         title,
		ContentLength: int(resp.ContentLength),
		Technologies:  technologies,
	}

	// Technology fingerprinting (when enabled via config)
	if cfg.TechDetect {
		result.DetectedTech = FingerprintTechnologies(resp, body)
	}

	return result, nil
}

// scanPorts performs port scanning on a single host
func scanPorts(ctx context.Context, host string, cfg *config.Config) (types.PortResult, error) {
	// Basic port scanning implementation
	// In a real implementation, you would use a proper port scanning library
	// For now, we'll just return a basic result structure

	commonPorts := []int{21, 22, 23, 25, 53, 80, 110, 143, 443, 993, 995, 8080, 8443}
	var openPorts []types.Port

	for _, port := range commonPorts {
		select {
		case <-ctx.Done():
			return types.PortResult{}, ctx.Err()
		default:
			// Basic port check (in real implementation, use proper TCP connection)
			if isPortOpen(host, port) {
				openPorts = append(openPorts, types.Port{
					Number:   port,
					Protocol: "tcp",
					State:    "open",
					Service:  getServiceName(port),
					Version:  "",
				})
			}
		}
	}

	return types.PortResult{
		Host:  host,
		IP:    "", // Would be resolved in real implementation
		Ports: openPorts,
	}, nil
}

// extractTitleFromBody extracts the title from an already-read HTTP response body.
func extractTitleFromBody(body []byte) string {
	if len(body) == 0 {
		return ""
	}
	re := regexp.MustCompile(`(?i)<title[^>]*>([^<]+)</title>`)
	matches := re.FindSubmatch(body)
	if len(matches) > 1 {
		return strings.TrimSpace(string(matches[1]))
	}
	return ""
}

// extractTechnologies extracts technologies from an HTTP response
func extractTechnologies(resp *http.Response) []string {
	var technologies []string

	// Check for common technology headers
	if server := resp.Header.Get("Server"); server != "" {
		technologies = append(technologies, server)
	}

	if poweredBy := resp.Header.Get("X-Powered-By"); poweredBy != "" {
		technologies = append(technologies, poweredBy)
	}

	return technologies
}

// isPortOpen checks if a port is open (basic implementation)
func isPortOpen(host string, port int) bool {
	// Basic TCP connection check
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), 3*time.Second)
	if err != nil {
		return false
	}
	defer func() { _ = conn.Close() }()
	return true
}

// enrichWithFingerprinting runs our fingerprint engine on HTTP results from
// the external httpx scanner. It makes a follow-up HTTP request per unique host
// to detect technologies from headers, cookies, and HTML body patterns.
func enrichWithFingerprinting(cfg *config.Config, results []types.HTTPResult) {
	if len(results) == 0 {
		return
	}

	client := &http.Client{
		Timeout: time.Duration(cfg.Timeout) * time.Second,
	}

	pool := utils.NewWorkerPool(cfg.Threads, cfg.RateLimit)
	defer pool.Stop()

	var mu sync.Mutex
	var wg sync.WaitGroup

	// Deduplicate by host to avoid redundant requests
	seen := make(map[string]bool)
	hostTechs := make(map[string][]types.Technology)

	for i := range results {
		host := ExtractHostFromURL(results[i].URL)
		if seen[host] {
			continue
		}
		seen[host] = true
		url := results[i].URL

		wg.Add(1)
		pool.Submit(func() {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
			defer cancel()

			req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
			if err != nil {
				return
			}
			req.Header.Set("User-Agent", "SubdomainX/1.0")

			resp, err := client.Do(req)
			if err != nil {
				return
			}
			body, _ := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
			resp.Body.Close()

			techs := FingerprintTechnologies(resp, body)
			if len(techs) > 0 {
				mu.Lock()
				hostTechs[host] = techs
				mu.Unlock()
			}
		})
	}

	wg.Wait()

	// Apply detected technologies back to results
	for i := range results {
		host := ExtractHostFromURL(results[i].URL)
		if techs, ok := hostTechs[host]; ok {
			results[i].DetectedTech = techs
		}
	}
}

// getServiceName returns the service name for a port
func getServiceName(port int) string {
	services := map[int]string{
		21:   "ftp",
		22:   "ssh",
		23:   "telnet",
		25:   "smtp",
		53:   "dns",
		80:   "http",
		110:  "pop3",
		143:  "imap",
		443:  "https",
		993:  "imaps",
		995:  "pop3s",
		8080: "http-proxy",
		8443: "https-alt",
	}

	if service, exists := services[port]; exists {
		return service
	}
	return "unknown"
}
