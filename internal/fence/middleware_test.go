package fence_test

import (
	"context"
	"errors"
	"testing"

	"vaultop/internal/fence"
)

func noopWrite(_ context.Context, _, _ string) error { return nil }

func TestGuardedWriter_ValidToken_Passes(t *testing.T) {
	m := fence.New()
	tok := m.Issue("db/pass", "leader")

	gw := fence.NewGuardedWriter(m, "db/pass", noopWrite)
	if err := gw.Write(context.Background(), tok, "db/pass", "secret"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGuardedWriter_StaleToken_ReturnsErrStaleFence(t *testing.T) {
	m := fence.New()
	old := m.Issue("db/pass", "leader-1")
	_ = m.Issue("db/pass", "leader-2")

	gw := fence.NewGuardedWriter(m, "db/pass", noopWrite)
	err := gw.Write(context.Background(), old, "db/pass", "value")

	if !errors.Is(err, fence.ErrStaleFence) {
		t.Fatalf("expected ErrStaleFence, got %v", err)
	}
}

func TestGuardedWriter_PropagatesInnerError(t *testing.T) {
	m := fence.New()
	tok := m.Issue("res", "h")

	sentinel := errors.New("backend unavailable")
	failing := func(_ context.Context, _, _ string) error { return sentinel }

	gw := fence.NewGuardedWriter(m, "res", failing)
	err := gw.Write(context.Background(), tok, "res", "v")

	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}
