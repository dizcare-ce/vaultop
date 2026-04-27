package fence_test

import (
	"errors"
	"testing"

	"vaultop/internal/fence"
)

func TestIssue_IncrementsSequence(t *testing.T) {
	m := fence.New()

	t1 := m.Issue("db/password", "node-1")
	t2 := m.Issue("db/password", "node-2")

	if t1.Seq != 1 {
		t.Fatalf("expected seq 1, got %d", t1.Seq)
	}
	if t2.Seq != 2 {
		t.Fatalf("expected seq 2, got %d", t2.Seq)
	}
	if t1.HolderID != "node-1" {
		t.Fatalf("expected holder node-1, got %s", t1.HolderID)
	}
}

func TestIssue_IndependentResources(t *testing.T) {
	m := fence.New()

	a := m.Issue("res/a", "h1")
	b := m.Issue("res/b", "h1")

	if a.Seq != 1 || b.Seq != 1 {
		t.Fatalf("resources should have independent sequences: a=%d b=%d", a.Seq, b.Seq)
	}
}

func TestCheck_CurrentToken_Passes(t *testing.T) {
	m := fence.New()
	tok := m.Issue("secret/key", "worker")

	if err := m.Check("secret/key", tok); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestCheck_StaleToken_ReturnsErrStaleFence(t *testing.T) {
	m := fence.New()
	old := m.Issue("secret/key", "worker-1")
	_ = m.Issue("secret/key", "worker-2") // supersedes old

	err := m.Check("secret/key", old)
	if !errors.Is(err, fence.ErrStaleFence) {
		t.Fatalf("expected ErrStaleFence, got %v", err)
	}
}

func TestCheck_UnknownResource_ReturnsError(t *testing.T) {
	m := fence.New()
	tok := fence.Token{Seq: 1, HolderID: "ghost"}

	if err := m.Check("no/such", tok); err == nil {
		t.Fatal("expected error for unknown resource")
	}
}

func TestCurrent_ReturnsMostRecent(t *testing.T) {
	m := fence.New()
	_ = m.Issue("r", "a")
	last := m.Issue("r", "b")

	got, ok := m.Current("r")
	if !ok {
		t.Fatal("expected token to exist")
	}
	if got.Seq != last.Seq {
		t.Fatalf("expected seq %d, got %d", last.Seq, got.Seq)
	}
}

func TestCurrent_Missing_ReturnsFalse(t *testing.T) {
	m := fence.New()
	_, ok := m.Current("nope")
	if ok {
		t.Fatal("expected false for missing resource")
	}
}
