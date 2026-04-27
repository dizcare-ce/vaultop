// Package presign provides time-limited, signed URLs or tokens that grant
// temporary access to a specific secret key without exposing credentials.
//
// A presigned reference encodes the target key, an expiry timestamp, and an
// HMAC-SHA256 signature derived from a shared secret.  The recipient can
// verify authenticity and freshness before acting on the reference.
package presign

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"
)

// ErrExpired is returned when a presigned token has passed its expiry time.
var ErrExpired = errors.New("presign: token has expired")

// ErrInvalid is returned when the token signature does not match or the
// format is malformed.
var ErrInvalid = errors.New("presign: token is invalid")

// Token represents a signed, time-limited reference to a secret key.
type Token struct {
	Key       string    // the secret key the token grants access to
	ExpiresAt time.Time // UTC expiry timestamp
	sig       string    // base64-encoded HMAC signature (unexported)
}

// Signer creates and verifies presigned tokens using a shared secret.
type Signer struct {
	secret []byte
	clock  func() time.Time
}

// New returns a Signer backed by the provided HMAC secret.
// The secret should be at least 32 bytes of high-entropy data.
func New(secret []byte) (*Signer, error) {
	if len(secret) < 16 {
		return nil, errors.New("presign: secret must be at least 16 bytes")
	}
	return &Signer{secret: secret, clock: time.Now}, nil
}

// Issue creates a signed token granting access to key for the given duration.
func (s *Signer) Issue(key string, ttl time.Duration) (string, error) {
	if key == "" {
		return "", errors.New("presign: key must not be empty")
	}
	if ttl <= 0 {
		return "", errors.New("presign: ttl must be positive")
	}

	expiry := s.clock().UTC().Add(ttl).Unix()
	payload := buildPayload(key, expiry)
	sig := s.sign(payload)

	// Format: base64(key) + "." + expiry + "." + sig
	encodedKey := base64.RawURLEncoding.EncodeToString([]byte(key))
	return fmt.Sprintf("%s.%d.%s", encodedKey, expiry, sig), nil
}

// Verify parses and validates a presigned token string.
// Returns the Token on success, or ErrExpired / ErrInvalid on failure.
func (s *Signer) Verify(raw string) (Token, error) {
	parts := strings.SplitN(raw, ".", 3)
	if len(parts) != 3 {
		return Token{}, ErrInvalid
	}

	keyBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return Token{}, ErrInvalid
	}

	var expiry int64
	if _, err := fmt.Sscanf(parts[1], "%d", &expiry); err != nil {
		return Token{}, ErrInvalid
	}

	payload := buildPayload(string(keyBytes), expiry)
	expectedSig := s.sign(payload)

	if !hmac.Equal([]byte(parts[2]), []byte(expectedSig)) {
		return Token{}, ErrInvalid
	}

	expiresAt := time.Unix(expiry, 0).UTC()
	if s.clock().UTC().After(expiresAt) {
		return Token{}, ErrExpired
	}

	return Token{
		Key:       string(keyBytes),
		ExpiresAt: expiresAt,
		sig:       parts[2],
	}, nil
}

// buildPayload returns the canonical string that is signed.
func buildPayload(key string, expiry int64) string {
	return fmt.Sprintf("%s:%d", key, expiry)
}

// sign computes a base64url-encoded HMAC-SHA256 of the payload.
func (s *Signer) sign(payload string) string {
	mac := hmac.New(sha256.New, s.secret)
	mac.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
