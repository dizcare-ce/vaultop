package metrics

import (
	"fmt"
	"time"
)

// RotationFunc is the signature of a secret rotation operation.
type RotationFunc func(key string) error

// Timed wraps fn, records success/failure in t, and returns elapsed duration.
func Timed(t *Tracker, name string, fn func() error) (time.Duration, error) {
	start := time.Now()
	err := fn()
	elapsed := time.Since(start)
	if err != nil {
		t.RecordFailure(name)
	} else {
		t.RecordSuccess(name)
	}
	return elapsed, err
}

// WrapRotation returns a RotationFunc that tracks metrics for each key.
func WrapRotation(t *Tracker, prefix string, inner RotationFunc) RotationFunc {
	return func(key string) error {
		name := fmt.Sprintf("%s/%s", prefix, key)
		_, err := Timed(t, name, func() error {
			return inner(key)
		})
		return err
	}
}
