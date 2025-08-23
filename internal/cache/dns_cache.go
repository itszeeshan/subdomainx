package cache

import (
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
