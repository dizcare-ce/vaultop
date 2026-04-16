package rotation_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultop/internal/audit"
	"github.com/yourusername/vaultop/internal/history"
	"github.com/yourusername/vaultop/internal/notification"
	"github.com/yourusername/vaultop/internal/rotation"
)

func buildOpts(t *testing.T, gen rotation.Generator, dryRun bool) rotation.RunOptions {
	t.Helper()
	h, err := history.Load(filepath.Join(t.TempDir(), "hist.json"))
	if err != nil {
		t.Fatalf("history.Load: %v", err)
	}
	return rotation.RunOptions{
		Provider:  newStubProvider(),
		Audit:     audit.NewLogger(os.Stderr),
		History:   h,
		Notifier:  notification.Noop{},
		Generator: gen,
		DryRun:    dryRun,
	}
}

func TestRun_AllSucceed(t *testing.T) {
	policies := []rotation.Policy{
		{Key: "alpha", Length: 16},
		{Key: "beta", Length: 16},
	}
	opts := buildOpts(t, rotation.FixedGenerator("s3cr3t"), false)
	res, err := rotation.Run(context.Background(), policies, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Rotated) != 2 {
		t.Errorf("expected 2 rotated, got %d", len(res.Rotated))
	}
	if len(res.Failed) != 0 {
		t.Errorf("expected no failures, got %v", res.Failed)
	}
}

func TestRun_GeneratorError_RecordsFailure(t *testing.T) {
	policies := []rotation.Policy{{Key: "gamma", Length: 16}}
	opts := buildOpts(t, rotation.ErrorGenerator(errors.New("gen fail")), false)
	res, err := rotation.Run(context.Background(), policies, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Failed) != 1 {
		t.Errorf("expected 1 failure, got %d", len(res.Failed))
	}
}

func TestRun_DryRun_DoesNotWriteHistory(t *testing.T) {
	policies := []rotation.Policy{{Key: "delta", Length: 16}}
	opts := buildOpts(t, rotation.FixedGenerator("x"), true)
	_, err := rotation.Run(context.Background(), policies, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, found := opts.History.LastRotated("delta")
	if found {
		t.Error("dry-run should not persist history")
	}
}

func TestRun_InvalidPolicy_RecordsFailure(t *testing.T) {
	policies := []rotation.Policy{{Key: "", Length: 16}}
	opts := buildOpts(t, rotation.FixedGenerator("x"), false)
	res, _ := rotation.Run(context.Background(), policies, opts)
	if len(res.Failed) != 1 {
		t.Errorf("expected 1 failure for invalid policy, got %d", len(res.Failed))
	}
}
