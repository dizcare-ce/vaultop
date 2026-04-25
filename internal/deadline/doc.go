// Package deadline provides per-key operation deadline tracking for vaultop.
//
// A Tracker registers a time-bounded window for a named operation (identified
// by a string key). Callers can check at any point whether the deadline has
// been exceeded, retrieve all active (in-progress) operations, or collect
// violations for alerting or audit purposes.
//
// Typical usage:
//
//	tr := deadline.New()
//	tr.Start("rotate/db-password", 30*time.Second)
//	defer tr.Finish("rotate/db-password")
//
//	if err := tr.Check("rotate/db-password"); err != nil {
//		// deadline exceeded — abort or alert
//	}
//
// Tracker is safe for concurrent use.
package deadline
