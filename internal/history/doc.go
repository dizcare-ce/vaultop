// Package history provides a lightweight, file-backed store for
// secret rotation history.
//
// Each rotation event — whether successful or not — is appended as a
// JSON entry to a local file. The store is safe for concurrent use
// within a single process.
//
// Typical usage:
//
//	s, err := history.Load("/var/lib/vaultop/history.json")
//	if err != nil { ... }
//
//	last := s.LastRotated("db/password")
//	// pass last to schedule.IsDueAt to decide whether rotation is needed
//
//	_ = s.Record(history.Entry{
//		SecretKey: "db/password",
//		RotatedAt: time.Now().UTC(),
//		Provider:  "aws",
//		Success:   true,
//	})
package history
