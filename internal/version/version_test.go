package version_test

import (
	"testing"

	"github.com/iamcalledned/vaultop/internal/version"
)

func TestRecord_IncrementsVersion(t *testing.T) {
	s := version.New()
	e1 := s.Record("db/pass", "alpha")
	e2 := s.Record("db/pass", "beta")

	if e1.Version != 1 {
		t.Fatalf("expected version 1, got %d", e1.Version)
	}
	if e2.Version != 2 {
		t.Fatalf("expected version 2, got %d", e2.Version)
	}
}

func TestGet_ReturnsCorrectEntry(t *testing.T) {
	s := version.New()
	s.Record("k", "v1")
	s.Record("k", "v2")

	e, err := s.Get("k", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Value != "v1" {
		t.Fatalf("expected v1, got %s", e.Value)
	}
}

func TestGet_UnknownKey_ReturnsErrNotFound(t *testing.T) {
	s := version.New()
	_, err := s.Get("missing", 1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGet_OutOfRange_ReturnsErrNotFound(t *testing.T) {
	s := version.New()
	s.Record("k", "only")
	_, err := s.Get("k", 99)
	if err == nil {
		t.Fatal("expected error for out-of-range version")
	}
}

func TestLatest_ReturnsMostRecent(t *testing.T) {
	s := version.New()
	s.Record("k", "first")
	s.Record("k", "second")
	s.Record("k", "third")

	e, err := s.Latest("k")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Value != "third" {
		t.Fatalf("expected third, got %s", e.Value)
	}
	if e.Version != 3 {
		t.Fatalf("expected version 3, got %d", e.Version)
	}
}

func TestLatest_UnknownKey_ReturnsErrNotFound(t *testing.T) {
	s := version.New()
	_, err := s.Latest("ghost")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestList_ReturnsAllInOrder(t *testing.T) {
	s := version.New()
	values := []string{"a", "b", "c"}
	for _, v := range values {
		s.Record("k", v)
	}

	entries := s.List("k")
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	for i, e := range entries {
		if e.Value != values[i] {
			t.Errorf("index %d: expected %s, got %s", i, values[i], e.Value)
		}
	}
}

func TestCount_ReturnsCorrectTotal(t *testing.T) {
	s := version.New()
	if s.Count("k") != 0 {
		t.Fatal("expected 0 for unknown key")
	}
	s.Record("k", "x")
	s.Record("k", "y")
	if s.Count("k") != 2 {
		t.Fatalf("expected 2, got %d", s.Count("k"))
	}
}
