// Package retry provides configurable retry logic for secret provider operations.
package retry

import (
	"errors"
	"fmt"
	"time"
)

// ErrMaxAttemptsReached is returned when all retry attempts are exhausted.
var ErrMaxAttemptsReached = errors.New("retry: max attempts reached")

// Config holds retry behaviour parameters.
type Config struct {
	MaxAttempts int
	Delay       time.Duration
	Multiplier  float64 // backoff multiplier; 1.0 = constant delay
}

// DefaultConfig returns a sensible default retry configuration.
func DefaultConfig() Config {
	return Config{
		MaxAttempts: 3,
		Delay:       200 * time.Millisecond,
		Multiplier:  2.0,
	}
}

// Do executes fn up to cfg.MaxAttempts times, applying exponential backoff.
// It returns nil on the first success, or ErrMaxAttemptsReached wrapping the
// last error after all attempts are exhausted.
func Do(cfg Config, fn func() error) error {
	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = 1
	}
	delay := cfg.Delay
	var last error
	for i := 0; i < cfg.MaxAttempts; i++ {
		if err := fn(); err == nil {
			return nil
		} else {
			last = err
		}
		if i < cfg.MaxAttempts-1 {
			time.Sleep(delay)
			if cfg.Multiplier > 0 {
				delay = time.Duration(float64(delay) * cfg.Multiplier)
			}
		}
	}
	return fmt.Errorf("%w: %v", ErrMaxAttemptsReached, last)
}
