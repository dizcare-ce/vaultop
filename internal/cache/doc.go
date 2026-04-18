// Package cache implements a lightweight TTL-based in-memory cache for secret
// values retrieved from cloud providers.
//
// During a single vaultop session multiple operations (validate, diff, export)
// may need to read the same secret key. The cache avoids redundant provider
// API calls by storing values for a configurable duration.
//
// Usage:
//
//	c := cache.New(30 * time.Second)
//	c.Set("myapp/db-password", plaintext)
//	if v, ok := c.Get("myapp/db-password"); ok {
//		// use cached value
//	}
//
// The cache is safe for concurrent use. Entries are considered expired after
// the TTL elapses; expired entries are treated as misses on Get but are not
// actively evicted — call Flush to clear all entries eagerly.
package cache
