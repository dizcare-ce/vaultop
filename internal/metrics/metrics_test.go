package metrics_test

import (
	"bytes"
	"strings"
	"testing"

	"vaultop/internal/metrics"
)

func TestRecordSuccess_Increments(t *testing.T) {
	tr := metrics.New()
	tr.RecordSuccess("rotate/db")
	tr.RecordSuccess("rotate/db")
	c, ok := tr.Get("rotate/db")
	if !ok {
		t.Fatal("expected counter to exist")
	}
	if c.Successes != 2 {
		t.Fatalf("want 2 successes, got %d", c.Successes)
	}
	if c.Failures != 0 {
		t.Fatalf("want 0 failures, got %d", c.Failures)
	}
}

func TestRecordFailure_Increments(t *testing.T) {
	tr := metrics.New()
	tr.RecordFailure("rotate/api")
	c, ok := tr.Get("rotate/api")
	if !ok {
		t.Fatal("expected counter to exist")
	}
	if c.Failures != 1 {
		t.Fatalf("want 1 failure, got %d", c.Failures)
	}
}

func TestGet_Missing(t *testing.T) {
	tr := metrics.New()
	_, ok := tr.Get("nonexistent")
	if ok {
		t.Fatal("expected not found")
	}
}

func TestAll_ReturnsAll(t *testing.T) {
	tr := metrics.New()
	tr.RecordSuccess("a")
	tr.RecordSuccess("b")
	tr.RecordFailure("b")
	all := tr.All()
	if len(all) != 2 {
		t.Fatalf("want 2 counters, got %d", len(all))
	}
}

func TestWriteSummary_ContainsName(t *testing.T) {
	tr := metrics.New()
	tr.RecordSuccess("rotate/db")
	var buf bytes.Buffer
	tr.WriteSummary(&buf)
	if !strings.Contains(buf.String(), "rotate/db") {
		t.Fatalf("summary missing key, got: %s", buf.String())
	}
}
