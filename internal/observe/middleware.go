package observe

import (
	"time"
)

// Timed wraps fn, recording the operation under key into obs.
// The original error (if any) is returned unchanged so callers can still
// inspect it.
func Timed(obs *Observer, key string, fn func() error) error {
	start := obs.clock()
	err := fn()
	obs.Record(key, time.Since(start), err)
	return err
}

// ProviderFunc is the signature used by provider Set/Get/Delete operations.
type ProviderFunc func(key, value string) error

// WrapSet returns a ProviderFunc that records each Set call into obs under
// the pattern "set:<key>".
func WrapSet(obs *Observer, fn ProviderFunc) ProviderFunc {
	return func(key, value string) error {
		return Timed(obs, "set:"+key, func() error {
			return fn(key, value)
		})
	}
}

// GetFunc is the signature for provider Get operations.
type GetFunc func(key string) (string, error)

// WrapGet returns a GetFunc that records each Get call into obs under
// the pattern "get:<key>".
func WrapGet(obs *Observer, fn GetFunc) GetFunc {
	return func(key string) (string, error) {
		var val string
		err := Timed(obs, "get:"+key, func() error {
			v, e := fn(key)
			val = v
			return e
		})
		return val, err
	}
}
