// Package diff provides utilities for comparing two snapshots of secrets.
//
// It detects which keys were added, removed, or had their values changed
// between a "before" and "after" state. This is useful for auditing
// changes after a rotation or a manual edit, and for producing human-readable
// change summaries in CLI output or audit logs.
//
// Usage:
//
//	changes := diff.Compare(before, after)
//	for _, c := range changes {
//		fmt.Println(c)
//	}
package diff
