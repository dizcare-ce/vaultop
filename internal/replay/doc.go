// Package replay provides idempotency protection for secret operations by
// detecting duplicate request IDs within a configurable time window.
//
// Usage:
//
//	detector := replay.New(5 * time.Minute)
//
//	if err := detector.Check(requestID); err == replay.ErrReplay {
//		// reject the request — already processed
//	}
//
// IDs are automatically evicted once the window has elapsed, keeping memory
// usage proportional to the request rate rather than total request count.
package replay
