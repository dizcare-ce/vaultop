package jitter_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/vaultop/internal/jitter"
)

func seeded() *rand.Rand {
	return rand.New(rand.NewSource(42))
}

func TestStrategy_IsValid(t *testing.T) {
	valid := []jitter.Strategy{jitter.Full, jitter.Equal, jitter.Decorrelated}
	for _, s := range valid {
		if !s.IsValid() {
			t.Errorf("expected strategy %d to be valid", s)
		}
	}
	if jitter.Strategy(99).IsValid() {
		t.Error("expected unknown strategy to be invalid")
	}
}

func TestApply_Full_WithinRange(t *testing.T) {
	cfg := jitter.DefaultConfig()
	base := 100 * time.Millisecond
	for i := 0; i < 200; i++ {
		d := jitter.Apply(cfg, base, 0, seeded())
		if d < 0 || d >= base {
			t.Fatalf("Full jitter %v out of [0, %v)", d, base)
		}
	}
}

func TestApply_Equal_WithinRange(t *testing.T) {
	cfg := jitter.Config{Strategy: jitter.Equal}
	base := 200 * time.Millisecond
	for i := 0; i < 200; i++ {
		d := jitter.Apply(cfg, base, 0, seeded())
		if d < base/2 || d > base {
			t.Fatalf("Equal jitter %v out of [%v, %v]", d, base/2, base)
		}
	}
}

func TestApply_Decorrelated_GreaterThanBase(t *testing.T) {
	cfg := jitter.Config{Strategy: jitter.Decorrelated}
	base := 50 * time.Millisecond
	last := base
	for i := 0; i < 50; i++ {
		d := jitter.Apply(cfg, base, last, seeded())
		if d < base {
			t.Fatalf("Decorrelated jitter %v should be >= base %v", d, base)
		}
		last = d
	}
}

func TestApply_Cap_IsRespected(t *testing.T) {
	cap := 30 * time.Millisecond
	cfg := jitter.Config{Strategy: jitter.Full, Cap: cap}
	base := 200 * time.Millisecond
	for i := 0; i < 200; i++ {
		d := jitter.Apply(cfg, base, 0, seeded())
		if d > cap {
			t.Fatalf("jitter %v exceeds cap %v", d, cap)
		}
	}
}

func TestApply_ZeroBase_ReturnsZero(t *testing.T) {
	cfg := jitter.DefaultConfig()
	d := jitter.Apply(cfg, 0, 0, seeded())
	if d != 0 {
		t.Fatalf("expected 0, got %v", d)
	}
}

func TestApply_NilRand_DoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panicked: %v", r)
		}
	}()
	cfg := jitter.DefaultConfig()
	_ = jitter.Apply(cfg, 100*time.Millisecond, 0, nil)
}
