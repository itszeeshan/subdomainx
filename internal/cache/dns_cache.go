package cache

import (
	"net"
	"sync"
	"time"
)

type DNSCache struct {
	cache map[string]cacheEntry
	mutex sync.RWMutex
}

type cacheEntry struct {
	ips       []string
	timestamp time.Time
}

func NewDNSCache() *DNSCache {
	return &DNSCache{
		cache: make(map[string]cacheEntry),
	}
}

func (d *DNSCache) Store(domain string, ips []string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.cache[domain] = cacheEntry{
		ips:       ips,
		timestamp: time.Now(),
	}
}

func (d *DNSCache) Lookup(domain string) []string {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if entry, exists := d.cache[domain]; exists {
		return entry.ips
	}

	return nil
}

// Resolve performs DNS resolution for a domain and caches the result
func (d *DNSCache) Resolve(domain string) []string {
	// Check cache first
	if ips := d.Lookup(domain); ips != nil {
		return ips
	}

	// Perform DNS resolution
	ips, err := net.LookupIP(domain)
	if err != nil {
		// Store empty result to avoid repeated failed lookups
		d.Store(domain, []string{})
		return []string{}
	}

	// Convert IP addresses to strings
	var ipStrings []string
	for _, ip := range ips {
		ipStrings = append(ipStrings, ip.String())
	}

	// Store in cache
	d.Store(domain, ipStrings)
	return ipStrings
}

func (d *DNSCache) Cleanup(maxAge time.Duration) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	now := time.Now()
	for domain, entry := range d.cache {
		if now.Sub(entry.timestamp) > maxAge {
			delete(d.cache, domain)
		}
	}
}
