package rotation_test

import (
	"context"
	"testing"

	"github.com/vaultop/internal/provider"
	"github.com/vaultop/internal/rotation"
)

func newStubProvider(t *testing.T) provider.Provider {
	t.Helper()
	p, err := provider.New(provider.TypeStub, nil)
	if err != nil {
		t.Fatalf("provider.New: %v", err)
	}
	return p
}

func TestRotate_SetsNewValues(t *testing.T) {
	p := newStubProvider(t)
	ctx := context.Background()

	r := rotation.New(p, rotation.Options{Generator: rotation.FixedGenerator("new-val")})
	results := r.Rotate(ctx, []string{"sec/a", "sec/b"})

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, res := range results {
		if res.Err != nil {
			t.Errorf("unexpected error for %s: %v", res.SecretID, res.Err)
		}
		got, err := p.Get(ctx, res.SecretID)
		if err != nil {
			t.Fatalf("Get %s: %v", res.SecretID, err)
		}
		if got != "new-val" {
			t.Errorf("expected 'new-val', got %q", got)
		}
	}
}

func TestRotate_DryRun_DoesNotPersist(t *testing.T) {
	p := newStubProvider(t)
	ctx := context.Background()
	_ = p.Set(ctx, "sec/x", "original")

	r := rotation.New(p, rotation.Options{DryRun: true, Generator: rotation.FixedGenerator("changed")})
	results := r.Rotate(ctx, []string{"sec/x"})

	if results[0].Err != nil {
		t.Fatalf("unexpected error: %v", results[0].Err)
	}
	got, _ := p.Get(ctx, "sec/x")
	if got != "original" {
		t.Errorf("dry-run should not change value, got %q", got)
	}
}

func TestRotate_GeneratorError_PropagatesResult(t *testing.T) {
	p := newStubProvider(t)
	r := rotation.New(p, rotation.Options{Generator: rotation.ErrorGenerator("boom")})
	results := r.Rotate(context.Background(), []string{"sec/fail"})

	if results[0].Err == nil {
		t.Error("expected error, got nil")
	}
}
