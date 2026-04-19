// Package retry implements a simple configurable retry mechanism with
// exponential backoff, intended for use when interacting with remote
// secret providers that may experience transient failures.
//
// Basic usage:
//
//	cfg := retry.DefaultConfig()
//	err := retry.Do(cfg, func() error {
//		return provider.Set(ctx, key, value)
//	})
package retry
