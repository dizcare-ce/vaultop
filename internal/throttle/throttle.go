// Package throttle provides a token-bucket throttler for limiting
// the rate at which secret operations are performed per caller.
package throttle

import (
	"fmt"
	"sync"
	"time"
)

// ErrThrottled is returned when the caller has exceeded their allowed burst.
var ErrThrottled = fmt.Errorf("throttle: request denied — rate limit exceeded")

// Config holds token-bucket parameters.
type Config struct {
	// Rate is the number of tokens replenished per second.
	Rate float64
	// Burst is the maximum number of tokens that can accumulate.
	Burst int
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{Rate: 5, Burst: 10}
}

type bucket struct {
	tokens   float64
	lastSeen time.Time
}

// Throttler enforces per-key token-bucket rate limiting.
type Throttler struct {
	cfg   Config
	mu    sync.Mutex
	state map[string]*bucket
	now   func() time.Time
}

// New creates a Throttler with the given config.
func New(cfg Config) *Throttler {
	return &Throttler{
		cfg:   cfg,
		state: make(map[string]*bucket),
		now:   time.Now,
	}
}

// Allow returns nil if the caller identified by key may proceed, or
// ErrThrottled if the token bucket is empty.
func (t *Throttler) Allow(key string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	b, ok := t.state[key]
	if !ok {
		b = &bucket{tokens: float64(t.cfg.Burst), lastSeen: now}
		t.state[key] = b
	}

	// Replenish tokens based on elapsed time.
	elapsed := now.Sub(b.lastSeen).Seconds()
	b.tokens += elapsed * t.cfg.Rate
	if b.tokens > float64(t.cfg.Burst) {
		b.tokens = float64(t.cfg.Burst)
	}
	b.lastSeen = now

	if b.tokens < 1 {
		return ErrThrottled
	}
	b.tokens--
	return nil
}

// Remaining returns the current token count for key (floored to zero).
func (t *Throttler) Remaining(key string) float64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	b, ok := t.state[key]
	if !ok {
		return float64(t.cfg.Burst)
	}
	return b.tokens
}

// Reset removes the bucket for key, restoring it to a full burst.
func (t *Throttler) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.state, key)
}
