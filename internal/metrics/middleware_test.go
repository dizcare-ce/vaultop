package metrics_test

import (
	"errors"
	"testing"

	"vaultop/internal/metrics"
)

func TestTimed_Success_RecordsSuccess(t *testing.T) {
	tr := metrics.New()
	_, err := metrics.Timed(tr, "op", func() error { return nil })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	c, _ := tr.Get("op")
	if c.Successes != 1 {
		t.Fatalf("want 1 success, got %d", c.Successes)
	}
}

func TestTimed_Failure_RecordsFailure(t *testing.T) {
	tr := metrics.New()
	sentinel := errors.New("boom")
	_, err := metrics.Timed(tr, "op", func() error { return sentinel })
	if !errors.Is(err, sentinel) {
		t.Fatalf("want sentinel error, got %v", err)
	}
	c, _ := tr.Get("op")
	if c.Failures != 1 {
		t.Fatalf("want 1 failure, got %d", c.Failures)
	}
}

func TestWrapRotation_TracksPerKey(t *testing.T) {
	tr := metrics.New()
	inner := func(key string) error {
		if key == "bad" {
			return errors.New("fail")
		}
		return nil
	}
	wrapped := metrics.WrapRotation(tr, "rotate", inner)
	_ = wrapped("good")
	_ = wrapped("bad")

	good, _ := tr.Get("rotate/good")
	if good.Successes != 1 {
		t.Fatalf("want 1 success for good, got %d", good.Successes)
	}
	bad, _ := tr.Get("rotate/bad")
	if bad.Failures != 1 {
		t.Fatalf("want 1 failure for bad, got %d", bad.Failures)
	}
}
