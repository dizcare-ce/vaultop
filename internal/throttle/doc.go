// Package throttle implements a per-key token-bucket rate limiter
// designed to protect secret-operation endpoints from bursts of
// automated requests.
//
// Each unique key (e.g. a caller identity or secret path) maintains
// its own independent bucket. Tokens accumulate at a configurable
// rate up to a maximum burst capacity. When the bucket is empty,
// Allow returns ErrThrottled.
//
// Example:
//
//	th := throttle.New(throttle.DefaultConfig())
//	if err := th.Allow(callerID); err != nil {
//		// reject or queue the request
//	}
package throttle
