// Package ttl provides expiry tracking for secrets managed by vaultop.
//
// A TTLMap stores per-key expiration times and allows callers to:
//   - Register a TTL duration for a secret key
//   - Check whether a specific key has expired
//   - List all keys that have passed their expiry time
//
// TTLMap is not safe for concurrent use; callers must synchronise
// access when sharing across goroutines.
package ttl
