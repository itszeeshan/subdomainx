package cache

import (
	"fmt"
	"testing"
	"time"
)

func TestNewDNSCache(t *testing.T) {
	cache := NewDNSCache()
	if cache == nil {
		t.Fatal("NewDNSCache returned nil")
	}

	if cache.cache == nil {
		t.Error("Cache map should be initialized")
	}
}

func TestDNSCacheStoreAndLookup(t *testing.T) {
	cache := NewDNSCache()
	domain := "example.com"
	ips := []string{"192.168.1.1", "10.0.0.1"}

	// Test storing
	cache.Store(domain, ips)

	// Test lookup
	result := cache.Lookup(domain)
	if result == nil {
		t.Fatal("Expected IPs, got nil")
	}

	if len(result) != len(ips) {
		t.Errorf("Expected %d IPs, got %d", len(ips), len(result))
	}

	for i, ip := range ips {
		if result[i] != ip {
			t.Errorf("Expected IP %d to be %s, got %s", i, ip, result[i])
		}
	}
}

func TestDNSCacheLookupNonExistent(t *testing.T) {
	cache := NewDNSCache()
	result := cache.Lookup("nonexistent.com")
	if result != nil {
		t.Errorf("Expected nil for non-existent domain, got %v", result)
	}
}

func TestDNSCacheStoreEmpty(t *testing.T) {
	cache := NewDNSCache()
	domain := "example.com"

	// Store empty slice
	cache.Store(domain, []string{})

	// Lookup should return empty slice, not nil
	result := cache.Lookup(domain)
	if result == nil {
		t.Error("Expected empty slice, got nil")
	}

	if len(result) != 0 {
		t.Errorf("Expected empty slice, got %d items", len(result))
	}
}

func TestDNSCacheResolve(t *testing.T) {
	cache := NewDNSCache()
	domain := "google.com" // Use a domain that should resolve

	// Test DNS resolution
	result := cache.Resolve(domain)
	if result == nil {
		t.Fatal("Expected IPs from DNS resolution, got nil")
	}

	// Should have at least one IP
	if len(result) == 0 {
		t.Error("Expected at least one IP from DNS resolution")
	}

	// Verify the result is cached
	cachedResult := cache.Lookup(domain)
	if cachedResult == nil {
		t.Error("Expected cached result, got nil")
	}

	if len(cachedResult) != len(result) {
		t.Errorf("Expected %d cached IPs, got %d", len(result), len(cachedResult))
	}
}

func TestDNSCacheResolveInvalidDomain(t *testing.T) {
	cache := NewDNSCache()
	domain := "invalid-domain-that-does-not-exist-12345.com"

	// Test DNS resolution of invalid domain
	result := cache.Resolve(domain)
	if result == nil {
		t.Fatal("Expected empty slice for invalid domain, got nil")
	}

	// Should return empty slice for invalid domain
	if len(result) != 0 {
		t.Errorf("Expected empty slice for invalid domain, got %d IPs", len(result))
	}

	// Verify the empty result is cached
	cachedResult := cache.Lookup(domain)
	if cachedResult == nil {
		t.Error("Expected cached empty result, got nil")
	}

	if len(cachedResult) != 0 {
		t.Errorf("Expected empty cached result, got %d IPs", len(cachedResult))
	}
}

func TestDNSCacheResolveCached(t *testing.T) {
	cache := NewDNSCache()
	domain := "example.com"
	ips := []string{"192.168.1.1", "10.0.0.1"}

	// Store in cache first
	cache.Store(domain, ips)

	// Resolve should return cached result
	result := cache.Resolve(domain)
	if result == nil {
		t.Fatal("Expected cached IPs, got nil")
	}

	if len(result) != len(ips) {
		t.Errorf("Expected %d cached IPs, got %d", len(ips), len(result))
	}

	for i, ip := range ips {
		if result[i] != ip {
			t.Errorf("Expected cached IP %d to be %s, got %s", i, ip, result[i])
		}
	}
}

func TestDNSCacheCleanup(t *testing.T) {
	cache := NewDNSCache()
	domain1 := "example1.com"
	domain2 := "example2.com"
	ips := []string{"192.168.1.1"}

	// Store entries
	cache.Store(domain1, ips)
	cache.Store(domain2, ips)

	// Verify both entries exist
	if cache.Lookup(domain1) == nil {
		t.Error("Expected domain1 to exist in cache")
	}
	if cache.Lookup(domain2) == nil {
		t.Error("Expected domain2 to exist in cache")
	}

	// Cleanup with very short max age (should remove all entries)
	cache.Cleanup(1 * time.Nanosecond)

	// Verify both entries are removed
	if cache.Lookup(domain1) != nil {
		t.Error("Expected domain1 to be removed from cache")
	}
	if cache.Lookup(domain2) != nil {
		t.Error("Expected domain2 to be removed from cache")
	}
}

func TestDNSCacheCleanupPartial(t *testing.T) {
	cache := NewDNSCache()
	domain1 := "example1.com"
	domain2 := "example2.com"
	ips := []string{"192.168.1.1"}

	// Store entries
	cache.Store(domain1, ips)

	// Wait a bit
	time.Sleep(10 * time.Millisecond)

	cache.Store(domain2, ips)

	// Cleanup with age that should only remove domain1
	cache.Cleanup(5 * time.Millisecond)

	// Verify domain1 is removed but domain2 remains
	if cache.Lookup(domain1) != nil {
		t.Error("Expected domain1 to be removed from cache")
	}
	if cache.Lookup(domain2) == nil {
		t.Error("Expected domain2 to remain in cache")
	}
}

func TestDNSCacheConcurrentAccess(t *testing.T) {
	cache := NewDNSCache()
	done := make(chan bool, 10)

	// Test concurrent store operations
	for i := 0; i < 5; i++ {
		go func(id int) {
			domain := fmt.Sprintf("example%d.com", id)
			ips := []string{fmt.Sprintf("192.168.1.%d", id)}
			cache.Store(domain, ips)
			done <- true
		}(i)
	}

	// Test concurrent lookup operations
	for i := 0; i < 5; i++ {
		go func(id int) {
			domain := fmt.Sprintf("example%d.com", id)
			cache.Lookup(domain)
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify no panic occurred and cache is still functional
	if cache.Lookup("example0.com") == nil {
		t.Error("Cache should still be functional after concurrent access")
	}
}

func TestDNSCacheResolveLocalhost(t *testing.T) {
	cache := NewDNSCache()
	domain := "localhost"

	// Test DNS resolution of localhost
	result := cache.Resolve(domain)
	if result == nil {
		t.Fatal("Expected IPs from localhost resolution, got nil")
	}

	// Should have at least one IP (usually 127.0.0.1 or ::1)
	if len(result) == 0 {
		t.Error("Expected at least one IP from localhost resolution")
	}

	// Verify the result is cached
	cachedResult := cache.Lookup(domain)
	if cachedResult == nil {
		t.Error("Expected cached result for localhost, got nil")
	}
}
