package alert

import (
	"testing"
	"time"
)

type stubSource struct {
	data map[string]time.Time
}

func (s *stubSource) Expiry(key string) (time.Time, bool) {
	v, ok := s.data[key]
	return v, ok
}

func TestCheckAll_FiresForExpiring(t *testing.T) {
	now := time.Now()
	src := &stubSource{data: map[string]time.Time{
		"key/a": now.Add(12 * time.Hour),
		"key/b": now.Add(30 * 24 * time.Hour),
	}}
	rules := []Rule{
		{Key: "key/a", WarnWithin: 7 * 24 * time.Hour, CritWithin: 24 * time.Hour},
		{Key: "key/b", WarnWithin: 7 * 24 * time.Hour, CritWithin: 24 * time.Hour},
	}
	var fired []Alert
	n := &collectNotifier{&fired}
	errs := CheckAll(rules, src, n, now)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(fired) != 1 || fired[0].Key != "key/a" {
		t.Fatalf("expected one alert for key/a, got %v", fired)
	}
}

func TestCheckAll_MissingKeySkipped(t *testing.T) {
	now := time.Now()
	src := &stubSource{data: map[string]time.Time{}}
	rules := []Rule{{Key: "missing", WarnWithin: time.Hour, CritWithin: time.Minute}}
	var fired []Alert
	n := &collectNotifier{&fired}
	CheckAll(rules, src, n, now)
	if len(fired) != 0 {
		t.Fatalf("expected no alerts, got %d", len(fired))
	}
}

func TestFilterFired_ReturnsOnlyFired(t *testing.T) {
	now := time.Now()
	src := &stubSource{data: map[string]time.Time{
		"x": now.Add(1 * time.Hour),
		"y": now.Add(30 * 24 * time.Hour),
	}}
	rules := []Rule{
		{Key: "x", WarnWithin: 7 * 24 * time.Hour, CritWithin: 2 * time.Hour},
		{Key: "y", WarnWithin: 7 * 24 * time.Hour, CritWithin: 2 * time.Hour},
	}
	alerts := FilterFired(rules, src, now)
	if len(alerts) != 1 || alerts[0].Key != "x" {
		t.Fatalf("expected alert for x only, got %v", alerts)
	}
}

type collectNotifier struct{ alerts *[]Alert }

func (c *collectNotifier) Notify(a Alert) error {
	*c.alerts = append(*c.alerts, a)
	return nil
}
