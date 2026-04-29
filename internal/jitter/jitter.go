// Package jitter provides utilities for adding randomised jitter to
// durations, which helps avoid thundering-herd problems when many
// goroutines retry or poll at the same interval.
package jitter

import (
	"math/rand"
	"time"
)

// Strategy controls how jitter is applied to a base duration.
type Strategy int

const (
	// Full replaces the duration with a random value in [0, base).
	Full Strategy = iota
	// Equal adds a random value in [-base/2, +base/2) to the base.
	Equal
	// Decorrelated picks a random value between base and 3*last,
	// suitable for exponential back-off chains.
	Decorrelated
)

// IsValid reports whether s is a recognised Strategy.
func (s Strategy) IsValid() bool {
	return s == Full || s == Equal || s == Decorrelated
}

// Config holds parameters for jitter calculation.
type Config struct {
	Strategy Strategy
	// Cap is an optional upper bound on the returned duration.
	// A zero value means no cap is applied.
	Cap time.Duration
}

// DefaultConfig returns a Config using the Full strategy with no cap.
func DefaultConfig() Config {
	return Config{Strategy: Full}
}

// Apply returns a jittered duration derived from base.
// For Decorrelated, last should be the previous jittered value;
// for other strategies it is ignored.
func Apply(cfg Config, base, last time.Duration, r *rand.Rand) time.Duration {
	if r == nil {
		r = rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec
	}

	var d time.Duration
	switch cfg.Strategy {
	case Equal:
		half := base / 2
		d = base - half + time.Duration(r.Int63n(int64(half)+1))
	case Decorrelated:
		if last <= 0 {
			last = base
		}
		max := 3 * last
		if max <= base {
			max = base + 1
		}
		d = base + time.Duration(r.Int63n(int64(max-base)))
	default: // Full
		if base <= 0 {
			return 0
		}
		d = time.Duration(r.Int63n(int64(base)))
	}

	if cfg.Cap > 0 && d > cfg.Cap {
		d = cfg.Cap
	}
	return d
}
