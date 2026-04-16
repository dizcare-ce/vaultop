package diff_test

import (
	"testing"

	"github.com/yourusername/vaultop/internal/diff"
)

func TestCompare_NoChanges(t *testing.T) {
	before := map[string]string{"foo": "bar", "baz": "qux"}
	after := map[string]string{"foo": "bar", "baz": "qux"}
	changes := diff.Compare(before, after)
	if len(changes) != 0 {
		t.Fatalf("expected no changes, got %d", len(changes))
	}
}

func TestCompare_Added(t *testing.T) {
	before := map[string]string{}
	after := map[string]string{"newkey": "newval"}
	changes := diff.Compare(before, after)
	if len(changes) != 1 || changes[0].Kind != diff.Added || changes[0].Key != "newkey" {
		t.Fatalf("expected one Added change, got %+v", changes)
	}
}

func TestCompare_Removed(t *testing.T) {
	before := map[string]string{"gone": "value"}
	after := map[string]string{}
	changes := diff.Compare(before, after)
	if len(changes) != 1 || changes[0].Kind != diff.Removed || changes[0].Key != "gone" {
		t.Fatalf("expected one Removed change, got %+v", changes)
	}
}

func TestCompare_Changed(t *testing.T) {
	before := map[string]string{"key": "old"}
	after := map[string]string{"key": "new"}
	changes := diff.Compare(before, after)
	if len(changes) != 1 || changes[0].Kind != diff.Changed || changes[0].Key != "key" {
		t.Fatalf("expected one Changed change, got %+v", changes)
	}
}

func TestHasChanges_True(t *testing.T) {
	if !diff.HasChanges(map[string]string{"a": "1"}, map[string]string{"a": "2"}) {
		t.Fatal("expected HasChanges to return true")
	}
}

func TestHasChanges_False(t *testing.T) {
	if diff.HasChanges(map[string]string{"a": "1"}, map[string]string{"a": "1"}) {
		t.Fatal("expected HasChanges to return false")
	}
}

func TestChange_String(t *testing.T) {
	c := diff.Change{Key: "mykey", Kind: diff.Added}
	if c.String() != "[added] mykey" {
		t.Fatalf("unexpected string: %s", c.String())
	}
}
