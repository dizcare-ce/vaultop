package observe

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time { return func() time.Time { return t } }

func TestRecord_IncrementsOps(t *testing.T) {
	o := New()
	o.Record("db/pass", 5*time.Millisecond, nil)
	o.Record("db/pass", 3*time.Millisecond, nil)
	e, ok := o.Get("db/pass")
	if !ok {
		t.Fatal("expected entry")
	}
	if e.Ops != 2 {
		t.Fatalf("want 2 ops, got %d", e.Ops)
	}
	if e.Errors != 0 {
		t.Fatalf("want 0 errors, got %d", e.Errors)
	}
}

func TestRecord_CountsErrors(t *testing.T) {
	o := New()
	o.Record("api/key", time.Millisecond, errors.New("boom"))
	o.Record("api/key", time.Millisecond, nil)
	e, _ := o.Get("api/key")
	if e.Errors != 1 {
		t.Fatalf("want 1 error, got %d", e.Errors)
	}
}

func TestAvgLatency(t *testing.T) {
	o := New()
	o.Record("k", 10*time.Millisecond, nil)
	o.Record("k", 20*time.Millisecond, nil)
	e, _ := o.Get("k")
	if e.AvgLatency() != 15*time.Millisecond {
		t.Fatalf("want 15ms, got %s", e.AvgLatency())
	}
}

func TestGet_Missing(t *testing.T) {
	o := New()
	_, ok := o.Get("no-such-key")
	if ok {
		t.Fatal("expected miss")
	}
}

func TestAll_ReturnsSnapshot(t *testing.T) {
	o := New()
	o.Record("a", time.Millisecond, nil)
	o.Record("b", time.Millisecond, nil)
	if len(o.All()) != 2 {
		t.Fatalf("want 2 entries, got %d", len(o.All()))
	}
}

func TestWriteSummary_ContainsKey(t *testing.T) {
	o := New()
	o.Record("secret/token", 2*time.Millisecond, nil)
	var buf bytes.Buffer
	o.WriteSummary(&buf)
	if !strings.Contains(buf.String(), "secret/token") {
		t.Fatalf("summary missing key: %s", buf.String())
	}
}

func TestTimed_RecordsSuccess(t *testing.T) {
	o := New()
	_ = Timed(o, "op", func() error { return nil })
	e, ok := o.Get("op")
	if !ok || e.Ops != 1 || e.Errors != 0 {
		t.Fatalf("unexpected entry: %+v", e)
	}
}

func TestTimed_RecordsFailure(t *testing.T) {
	o := New()
	_ = Timed(o, "op", func() error { return errors.New("fail") })
	e, _ := o.Get("op")
	if e.Errors != 1 {
		t.Fatalf("want 1 error, got %d", e.Errors)
	}
}

func TestWrapGet_TracksPerKey(t *testing.T) {
	o := New()
	store := map[string]string{"x": "val"}
	get := WrapGet(o, func(k string) (string, error) {
		v, ok := store[k]
		if !ok {
			return "", errors.New("not found")
		}
		return v, nil
	})
	get("x") //nolint
	get("missing") //nolint
	if e, _ := o.Get("get:x"); e.Ops != 1 || e.Errors != 0 {
		t.Fatalf("unexpected: %+v", e)
	}
	if e, _ := o.Get("get:missing"); e.Errors != 1 {
		t.Fatalf("expected error recorded")
	}
}
