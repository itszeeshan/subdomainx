package scanner

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/itszeeshan/subdomainx/internal/config"
	"github.com/itszeeshan/subdomainx/internal/types"
	"github.com/itszeeshan/subdomainx/internal/utils"
)

// serviceFingerprint defines a known vulnerable service for takeover detection.
type serviceFingerprint struct {
	Name         string
	CnamePattern []string // suffix patterns to match CNAME targets
	BodyPattern  string   // string to look for in HTTP response body
	Risk         string   // "high", "medium", "low"
}

var serviceFingerprints = []serviceFingerprint{
	{
		Name:         "AWS S3",
		CnamePattern: []string{".s3.amazonaws.com", ".s3-website"},
		BodyPattern:  "NoSuchBucket",
		Risk:         "high",
	},
	{
		Name:         "GitHub Pages",
		CnamePattern: []string{".github.io"},
		BodyPattern:  "There isn't a GitHub Pages site here",
		Risk:         "high",
	},
	{
		Name:         "Heroku",
		CnamePattern: []string{".herokuapp.com", ".herokussl.com"},
		BodyPattern:  "No such app",
		Risk:         "high",
	},
	{
		Name:         "Azure",
		CnamePattern: []string{".azurewebsites.net", ".cloudapp.net", ".azure-api.net", ".azurehdinsight.net", ".azureedge.net", ".trafficmanager.net"},
		BodyPattern:  "",
		Risk:         "high",
	},
	{
		Name:         "Shopify",
		CnamePattern: []string{".myshopify.com"},
		BodyPattern:  "Sorry, this shop is currently unavailable",
		Risk:         "high",
	},
	{
		Name:         "Fastly",
		CnamePattern: []string{".fastly.net"},
		BodyPattern:  "Fastly error: unknown domain",
		Risk:         "high",
	},
	{
		Name:         "Pantheon",
		CnamePattern: []string{".pantheonsite.io"},
		BodyPattern:  "404 error unknown site",
		Risk:         "high",
	},
	{
		Name:         "Tumblr",
		CnamePattern: []string{".tumblr.com"},
		BodyPattern:  "There's nothing here",
		Risk:         "medium",
	},
	{
		Name:         "WordPress.com",
		CnamePattern: []string{".wordpress.com"},
		BodyPattern:  "Do you want to register",
		Risk:         "medium",
	},
	{
		Name:         "Fly.io",
		CnamePattern: []string{".fly.dev"},
		BodyPattern:  "",
		Risk:         "medium",
	},
	{
		Name:         "Surge.sh",
		CnamePattern: []string{".surge.sh"},
		BodyPattern:  "project not found",
		Risk:         "high",
	},
	{
		Name:         "Netlify",
		CnamePattern: []string{".netlify.app", ".netlify.com"},
		BodyPattern:  "Not Found - Request ID",
		Risk:         "high",
	},
}

// RunTakeoverCheck checks subdomains for potential takeover vulnerabilities
// by examining CNAME records and optionally HTTP responses.
func RunTakeoverCheck(cfg *config.Config, subdomains []types.SubdomainResult, httpResults []types.HTTPResult) ([]types.TakeoverResult, error) {
	if len(subdomains) == 0 {
		return nil, nil
	}

	// Build HTTP body lookup map from existing HTTP results
	httpBodyMap := make(map[string]string) // subdomain -> response body
	if len(httpResults) > 0 {
		httpBodyMap = fetchBodiesForTakeover(cfg, httpResults)
	}

	pool := utils.NewWorkerPool(cfg.Threads, cfg.RateLimit)
	defer pool.Stop()

	results := make(chan types.TakeoverResult, len(subdomains))
	var wg sync.WaitGroup

	for _, sub := range subdomains {
		sub := sub
		wg.Add(1)
		pool.Submit(func() {
			defer wg.Done()
			if result, ok := checkTakeover(sub.Subdomain, httpBodyMap); ok {
				results <- result
			}
		})
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var takeoverResults []types.TakeoverResult
	for result := range results {
		takeoverResults = append(takeoverResults, result)
	}

	return takeoverResults, nil
}

// checkTakeover checks a single subdomain for takeover vulnerability.
func checkTakeover(subdomain string, httpBodyMap map[string]string) (types.TakeoverResult, bool) {
	// Look up CNAME
	cname, err := net.LookupCNAME(subdomain)
	if err != nil || cname == "" || cname == subdomain+"." {
		return types.TakeoverResult{}, false
	}
	cname = strings.TrimSuffix(cname, ".")

	// Check if CNAME matches a known vulnerable service
	for _, fp := range serviceFingerprints {
		if !MatchesCnamePattern(cname, fp.CnamePattern) {
			continue
		}

		// Check if CNAME target is dangling (NXDOMAIN)
		isDangling := isCnameDangling(cname)

		// Check HTTP body if available
		bodyMatch := false
		if fp.BodyPattern != "" {
			if body, exists := httpBodyMap[subdomain]; exists {
				bodyMatch = strings.Contains(body, fp.BodyPattern)
			}
		}

		if isDangling || bodyMatch {
			evidence := buildEvidence(isDangling, bodyMatch, cname, fp)
			return types.TakeoverResult{
				Subdomain: subdomain,
				CNAME:     cname,
				Risk:      fp.Risk,
				Service:   fp.Name,
				Evidence:  evidence,
			}, true
		}
	}

	return types.TakeoverResult{}, false
}

// MatchesCnamePattern checks if a CNAME target matches any of the service patterns.
func MatchesCnamePattern(cname string, patterns []string) bool {
	cnameLower := strings.ToLower(cname)
	for _, pattern := range patterns {
		if strings.HasSuffix(cnameLower, strings.ToLower(pattern)) {
			return true
		}
	}
	return false
}

// isCnameDangling checks if a CNAME target fails to resolve (NXDOMAIN).
func isCnameDangling(cname string) bool {
	_, err := net.LookupHost(cname)
	if err != nil {
		if dnsErr, ok := err.(*net.DNSError); ok {
			return dnsErr.IsNotFound
		}
		// Treat other errors (timeout etc.) as not dangling
		return false
	}
	return false
}

// buildEvidence creates a human-readable evidence string.
func buildEvidence(isDangling, bodyMatch bool, cname string, fp serviceFingerprint) string {
	var parts []string
	if isDangling {
		parts = append(parts, fmt.Sprintf("CNAME %s → NXDOMAIN", cname))
	}
	if bodyMatch {
		parts = append(parts, fmt.Sprintf("HTTP response contains: %q", fp.BodyPattern))
	}
	return strings.Join(parts, "; ")
}

// fetchBodiesForTakeover fetches HTTP response bodies for the subdomains
// that have HTTP results, to check for takeover fingerprints in the body.
func fetchBodiesForTakeover(cfg *config.Config, httpResults []types.HTTPResult) map[string]string {
	bodies := make(map[string]string)
	client := &http.Client{
		Timeout: time.Duration(cfg.Timeout) * time.Second,
	}

	for _, hr := range httpResults {
		host := ExtractHostFromURL(hr.URL)
		if _, exists := bodies[host]; exists {
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
		req, err := http.NewRequestWithContext(ctx, "GET", hr.URL, nil)
		if err != nil {
			cancel()
			continue
		}
		req.Header.Set("User-Agent", "SubdomainX/1.0")

		resp, err := client.Do(req)
		if err != nil {
			cancel()
			continue
		}

		body, err := io.ReadAll(io.LimitReader(resp.Body, 32*1024))
		resp.Body.Close()
		cancel()
		if err != nil {
			continue
		}

		bodies[host] = string(body)
	}

	return bodies
}

// ExtractHostFromURL extracts the hostname from a URL string.
func ExtractHostFromURL(rawURL string) string {
	host := rawURL
	for _, prefix := range []string{"https://", "http://"} {
		host = strings.TrimPrefix(host, prefix)
	}
	host = strings.TrimRight(host, "/")
	if idx := strings.Index(host, ":"); idx != -1 {
		host = host[:idx]
	}
	if idx := strings.Index(host, "/"); idx != -1 {
		host = host[:idx]
	}
	return host
}

// PrintTakeoverSummary prints a formatted summary of takeover results to the terminal.
func PrintTakeoverSummary(results []types.TakeoverResult) {
	if len(results) == 0 {
		return
	}

	log.Println("")
	log.Println("⚠️  TAKEOVER RISK DETECTED:")
	log.Println("")

	highCount, medCount, lowCount := 0, 0, 0
	for _, r := range results {
		var icon string
		switch r.Risk {
		case "high":
			icon = "🔴"
			highCount++
		case "medium":
			icon = "🟠"
			medCount++
		default:
			icon = "🟡"
			lowCount++
		}
		log.Printf("  %s [%s] %s", icon, strings.ToUpper(r.Risk), r.Subdomain)
		if r.CNAME != "" {
			log.Printf("           CNAME → %s", r.CNAME)
		}
		log.Printf("           Service: %s — %s", r.Service, r.Evidence)
		log.Println()
	}

	log.Printf("Summary: %d high risk, %d medium risk, %d low risk\n", highCount, medCount, lowCount)
}
