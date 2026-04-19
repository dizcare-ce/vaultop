// Package quota provides per-key operation rate limiting using a fixed
// rolling window strategy. It is used to prevent excessive secret reads,
// writes, or rotations within a configured time period.
//
// Usage:
//
//	limiter := quota.New(quota.Config{MaxOps: 10, Window: time.Minute})
//	if err := limiter.Allow(secretKey); err != nil {
//		// handle quota exceeded
//	}
package quota
