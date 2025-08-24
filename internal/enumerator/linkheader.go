package enumerator

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/itszeeshan/subdomainx/internal/config"
	"github.com/tomnomnom/linkheader"
)

type LinkHeaderEnumerator struct {
	client *http.Client
}

func (l *LinkHeaderEnumerator) Name() string {
	return "linkheader"
}

func (l *LinkHeaderEnumerator) Enumerate(ctx context.Context, domain string, cfg *config.Config, ) ([]string, error) {
	// This enumerator works by checking discovered subdomains for Link headers
	// It needs to be run after other enumerators have discovered initial subdomains
	// For now, we'll check the main domain and common subdomains

	subdomainSet := make(map[string]bool)

	// Check main domain and common subdomains for Link headers
	targets := []string{
		domain,
		"www." + domain,
		"api." + domain,
		"docs." + domain,
		"dev." + domain,
		"staging." + domain,
		"test." + domain,
	}

	for _, target := range targets {
		subdomains, err := l.checkLinkHeaders(ctx, target, domain)
		if err != nil {
			// Continue with other targets even if one fails
			continue
		}

		for _, subdomain := range subdomains {
			subdomainSet[subdomain] = true
		}
	}

	// Convert set to slice
	var subdomains []string
	for subdomain := range subdomainSet {
		subdomains = append(subdomains, subdomain)
	}

	return subdomains, nil
}

func (l *LinkHeaderEnumerator) checkLinkHeaders(ctx context.Context, target, baseDomain string) ([]string, error) {
	var subdomains []string

	// Try HTTP first
	httpURL := fmt.Sprintf("http://%s", target)
	subdomains = append(subdomains, l.fetchLinkHeaders(ctx, httpURL, baseDomain)...)

	// Try HTTPS
	httpsURL := fmt.Sprintf("https://%s", target)
	subdomains = append(subdomains, l.fetchLinkHeaders(ctx, httpsURL, baseDomain)...)

	return subdomains, nil
}

func (l *LinkHeaderEnumerator) fetchLinkHeaders(ctx context.Context, urlStr, baseDomain string) []string {
	var subdomains []string

	req, err := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
	if err != nil {
		return subdomains
	}

	// Set reasonable headers
	req.Header.Set("User-Agent", "SubdomainX/1.2.0")
	req.Header.Set("Accept", "*/*")

	resp, err := l.client.Do(req)
	if err != nil {
		return subdomains
	}
	defer resp.Body.Close()

	// Parse Link headers
	linkHeader := resp.Header.Get("Link")
	if linkHeader == "" {
		return subdomains
	}

	links := linkheader.Parse(linkHeader)
	for _, link := range links {
		// Extract subdomains from Link URLs
		subdomain := l.extractSubdomainFromURL(link.URL, baseDomain)
		if subdomain != "" {
			subdomains = append(subdomains, subdomain)
		}
	}

	return subdomains
}

func (l *LinkHeaderEnumerator) extractSubdomainFromURL(urlStr, baseDomain string) string {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return ""
	}

	hostname := parsedURL.Hostname()
	if hostname == "" {
		return ""
	}

	// Check if it's a subdomain of our target domain
	if strings.HasSuffix(hostname, "."+baseDomain) && hostname != baseDomain {
		return hostname
	}

	return ""
}

func init() {
	// Create Link Header enumerator with default settings
	enumerator := &LinkHeaderEnumerator{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
	RegisterEnumerator(enumerator)
}
