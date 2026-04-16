package rotation

import (
	"errors"
	"time"
)

// Policy defines rotation behaviour for a single secret.
type Policy struct {
	// Key is the provider-specific secret identifier.
	Key string

	// Interval is how often the secret should be rotated.
	// A zero value disables automatic rotation.
	Interval time.Duration

	// Length is the desired byte-length of the generated secret.
	// Defaults to DefaultSecretLength when zero.
	Length int
}

// Validate returns an error if the policy is not usable.
func (p Policy) Validate() error {
	if p.Key == "" {
		return errors.New("rotation policy: key must not be empty")
	}
	if p.Length < 0 {
		return errors.New("rotation policy: length must be non-negative")
	}
	return nil
}

// EffectiveLength returns the length to use, falling back to the default.
func (p Policy) EffectiveLength() int {
	if p.Length == 0 {
		return DefaultSecretLength
	}
	return p.Length
}
