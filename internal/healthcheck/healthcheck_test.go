package healthcheck_test

import (
	"errors"
	"testing"

	"github.com/vaultop/internal/healthcheck"
)

func TestRunAll_SingleOK(t *testing.T) {
	c := healthcheck.New()
	c.Register("db", func() healthcheck.Status {
		return healthcheck.OK("db")
	})

	results, overall := c.RunAll()
	if overall != healthcheck.LevelOK {
		t.Fatalf("expected ok, got %s", overall)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Level != healthcheck.LevelOK {
		t.Errorf("expected ok level")
	}
}

func TestRunAll_UnhealthyDominates(t *testing.T) {
	c := healthcheck.New()
	c.Register("cache", func() healthcheck.Status {
		return healthcheck.OK("cache")
	})
	c.Register("vault", func() healthcheck.Status {
		return healthcheck.Unhealthy("vault", errors.New("connection refused"))
	})

	_, overall := c.RunAll()
	if overall != healthcheck.LevelUnhealthy {
		t.Fatalf("expected unhealthy, got %s", overall)
	}
}

func TestRunAll_DegradedBelowUnhealthy(t *testing.T) {
	c := healthcheck.New()
	c.Register("p1", func() healthcheck.Status {
		return healthcheck.Status{Level: healthcheck.LevelDegraded, Message: "slow"}
	})
	c.Register("p2", func() healthcheck.Status {
		return healthcheck.OK("p2")
	})

	_, overall := c.RunAll()
	if overall != healthcheck.LevelDegraded {
		t.Fatalf("expected degraded, got %s", overall)
	}
}

func TestRunAll_Empty(t *testing.T) {
	c := healthcheck.New()
	results, overall := c.RunAll()
	if overall != healthcheck.LevelOK {
		t.Fatalf("expected ok for empty checker, got %s", overall)
	}
	if len(results) != 0 {
		t.Fatalf("expected no results, got %d", len(results))
	}
}

func TestRegister_OverwritesPrevious(t *testing.T) {
	c := healthcheck.New()
	c.Register("svc", func() healthcheck.Status {
		return healthcheck.Unhealthy("svc", errors.New("first"))
	})
	c.Register("svc", func() healthcheck.Status {
		return healthcheck.OK("svc")
	})

	_, overall := c.RunAll()
	if overall != healthcheck.LevelOK {
		t.Fatalf("expected ok after overwrite, got %s", overall)
	}
}

func TestRunAll_LatencyPopulated(t *testing.T) {
	c := healthcheck.New()
	c.Register("fast", func() healthcheck.Status {
		return healthcheck.OK("fast")
	})

	results, _ := c.RunAll()
	if results[0].Latency < 0 {
		t.Errorf("latency should be non-negative")
	}
}
