package circuit_test

import (
	"testing"
	"time"

	"github.com/vaultop/internal/circuit"
)

func cfg(maxFail int, openDur time.Duration) circuit.Config {
	return circuit.Config{MaxFailures: maxFail, OpenDuration: openDur}
}

func TestAllow_ClosedState_Permits(t *testing.T) {
	b := circuit.New(cfg(3, time.Minute))
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestRecordFailure_OpensAfterThreshold(t *testing.T) {
	b := circuit.New(cfg(3, time.Minute))
	for i := 0; i < 3; i++ {
		b.RecordFailure()
	}
	if b.State() != circuit.StateOpen {
		t.Fatalf("expected open, got %s", b.State())
	}
	if err := b.Allow(); err != circuit.ErrOpen {
		t.Fatalf("expected ErrOpen, got %v", err)
	}
}

func TestRecordSuccess_ClosesCircuit(t *testing.T) {
	b := circuit.New(cfg(2, time.Minute))
	b.RecordFailure()
	b.RecordFailure()
	if b.State() != circuit.StateOpen {
		t.Fatal("expected open")
	}
	b.RecordSuccess()
	if b.State() != circuit.StateClosed {
		t.Fatalf("expected closed, got %s", b.State())
	}
}

func TestAllow_TransitionsToHalfOpen_AfterDuration(t *testing.T) {
	b := circuit.New(cfg(1, 10*time.Millisecond))
	b.RecordFailure()
	if b.State() != circuit.StateOpen {
		t.Fatal("expected open")
	}
	time.Sleep(20 * time.Millisecond)
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil after timeout, got %v", err)
	}
	if b.State() != circuit.StateHalfOpen {
		t.Fatalf("expected half-open, got %s", b.State())
	}
}

func TestState_String(t *testing.T) {
	cases := []struct {
		s    circuit.State
		want string
	}{
		{circuit.StateClosed, "closed"},
		{circuit.StateOpen, "open"},
		{circuit.StateHalfOpen, "half-open"},
	}
	for _, tc := range cases {
		if got := tc.s.String(); got != tc.want {
			t.Errorf("State(%d).String() = %q, want %q", tc.s, got, tc.want)
		}
	}
}

func TestDefaultConfig_SensibleValues(t *testing.T) {
	c := circuit.DefaultConfig()
	if c.MaxFailures <= 0 {
		t.Error("MaxFailures should be positive")
	}
	if c.OpenDuration <= 0 {
		t.Error("OpenDuration should be positive")
	}
}
