// Package circuit implements a simple circuit breaker for protecting
// provider calls from cascading failures.
package circuit

import (
	"errors"
	"sync"
	"time"
)

// ErrOpen is returned when the circuit breaker is in the open state.
var ErrOpen = errors.New("circuit breaker is open")

// State represents the current state of the circuit breaker.
type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

// String returns a human-readable state name.
func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// Config holds tuning parameters for the circuit breaker.
type Config struct {
	MaxFailures  int
	OpenDuration time.Duration
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		MaxFailures:  5,
		OpenDuration: 30 * time.Second,
	}
}

// Breaker is a thread-safe circuit breaker.
type Breaker struct {
	mu       sync.Mutex
	cfg      Config
	state    State
	failures int
	openedAt time.Time
}

// New creates a Breaker with the given config.
func New(cfg Config) *Breaker {
	return &Breaker{cfg: cfg}
}

// Allow reports whether a call should be attempted.
// It transitions from open to half-open once OpenDuration has elapsed.
func (b *Breaker) Allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case StateOpen:
		if time.Since(b.openedAt) >= b.cfg.OpenDuration {
			b.state = StateHalfOpen
			return nil
		}
		return ErrOpen
	}
	return nil
}

// RecordSuccess resets failure count and closes the circuit.
func (b *Breaker) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
	b.state = StateClosed
}

// RecordFailure increments the failure counter and may open the circuit.
func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures++
	if b.failures >= b.cfg.MaxFailures {
		b.state = StateOpen
		b.openedAt = time.Now()
	}
}

// State returns the current circuit state.
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}
