// Package diff provides utilities for comparing secret snapshots.
package diff

import "fmt"

// ChangeKind describes the type of change detected between two snapshots.
type ChangeKind string

const (
	Added   ChangeKind = "added"
	Removed ChangeKind = "removed"
	Changed ChangeKind = "changed"
)

// Change represents a single detected difference between two secret maps.
type Change struct {
	Key  string
	Kind ChangeKind
}

func (c Change) String() string {
	return fmt.Sprintf("[%s] %s", c.Kind, c.Key)
}

// Compare returns the list of changes between the before and after secret maps.
// Values are compared by equality; keys present only in before are Removed,
// only in after are Added, and present in both with different values are Changed.
func Compare(before, after map[string]string) []Change {
	var changes []Change

	for k, vBefore := range before {
		if vAfter, ok := after[k]; !ok {
			changes = append(changes, Change{Key: k, Kind: Removed})
		} else if vBefore != vAfter {
			changes = append(changes, Change{Key: k, Kind: Changed})
		}
	}

	for k := range after {
		if _, ok := before[k]; !ok {
			changes = append(changes, Change{Key: k, Kind: Added})
		}
	}

	return changes
}

// HasChanges returns true if any differences exist between before and after.
func HasChanges(before, after map[string]string) bool {
	return len(Compare(before, after)) > 0
}

// FilterByKind returns only the changes matching the given ChangeKind.
func FilterByKind(changes []Change, kind ChangeKind) []Change {
	var filtered []Change
	for _, c := range changes {
		if c.Kind == kind {
			filtered = append(filtered, c)
		}
	}
	return filtered
}
