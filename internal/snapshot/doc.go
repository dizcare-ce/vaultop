// Package snapshot provides point-in-time capture and restoration of secrets
// managed by a vaultop provider.
//
// A snapshot records the values of a specified set of secret keys at a given
// moment. Snapshots can be persisted to disk as JSON files and later loaded
// to restore secrets to a previous state — useful for rollback after a failed
// rotation or for auditing purposes.
//
// Basic usage:
//
//	s, err := snapshot.Take(p, []string{"db/password", "api/key"})
//	if err != nil { ... }
//	_ = snapshot.Save(s, "backup.json")
//
//	// later:
//	loaded, _ := snapshot.Load("backup.json")
//	_ = snapshot.Restore(p, loaded)
package snapshot
