// Package watermark embeds a tamper-evident marker into secret values so that
// out-of-band modifications can be detected when the secret is next read.
package watermark

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

const separator = "."

// ErrInvalidWatermark is returned when a value's embedded marker does not
// match the expected HMAC, indicating potential tampering.
var ErrInvalidWatermark = errors.New("watermark: invalid or missing marker")

// Manager signs and verifies watermarked secret values.
type Manager struct {
	secret []byte
}

// New returns a Manager that uses secret as the HMAC key.
// secret must be at least 16 bytes.
func New(secret []byte) (*Manager, error) {
	if len(secret) < 16 {
		return nil, errors.New("watermark: secret must be at least 16 bytes")
	}
	return &Manager{secret: secret}, nil
}

// Apply appends a base64-encoded HMAC tag to value and returns the marked string.
func (m *Manager) Apply(key, value string) string {
	tag := m.sign(key, value)
	return value + separator + tag
}

// Verify checks that marked was produced by Apply for the given key.
// It returns the original plain value on success.
func (m *Manager) Verify(key, marked string) (string, error) {
	idx := strings.LastIndex(marked, separator)
	if idx < 0 {
		return "", ErrInvalidWatermark
	}
	value := marked[:idx]
	tag := marked[idx+1:]
	expected := m.sign(key, value)
	if !hmac.Equal([]byte(tag), []byte(expected)) {
		return "", ErrInvalidWatermark
	}
	return value, nil
}

// IsMarked reports whether marked appears to contain a watermark tag.
func IsMarked(marked string) bool {
	return strings.Contains(marked, separator)
}

func (m *Manager) sign(key, value string) string {
	h := hmac.New(sha256.New, m.secret)
	_, _ = fmt.Fprintf(h, "%s:%s", key, value)
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}
