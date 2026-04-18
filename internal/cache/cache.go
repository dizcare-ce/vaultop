// Package cache provides a short-lived in-memory cache for secret values,
// reducing redundant provider round-trips during a single vaultop session.
package cache

import (
	"sync"
	"time"
)

// Entry holds a cached secret value and its expiry time.
type Entry struct {
	Value     string
	ExpiresAt time.Time
}

// Cache is a thread-safe in-memory store with TTL-based expiry.
type Cache struct {
	mu      sync.RWMutex
	items   map[string]Entry
	ttl     time.Duration
}

// New returns a Cache that expires entries after ttl.
func New(ttl time.Duration) *Cache {
	return &Cache{
		items: make(map[string]Entry),
		ttl:   ttl,
	}
}

// Set stores value under key, overwriting any existing entry.
func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = Entry{
		Value:     value,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

// Get returns the cached value for key and whether it was found and still valid.
func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.items[key]
	if !ok || time.Now().After(e.ExpiresAt) {
		return "", false
	}
	return e.Value, true
}

// Delete removes the entry for key, if present.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Flush removes all entries from the cache.
func (c *Cache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]Entry)
}

// Len returns the number of entries currently in the cache (including expired).
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}
