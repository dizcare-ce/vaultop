// Package token provides short-lived token generation and validation
// for service-to-service authentication when accessing secrets.
package token

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// ErrExpired is returned when a token has passed its TTL.
var ErrExpired = errors.New("token: expired")

// ErrInvalid is returned when a token signature does not match.
var ErrInvalid = errors.New("token: invalid signature")

// Token holds an opaque credential and its expiry.
type Token struct {
	Value     string
	ExpiresAt time.Time
}

// IsExpired reports whether the token is past its expiry time.
func (t Token) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// Manager issues and validates HMAC-signed tokens.
type Manager struct {
	secret []byte
	ttl    time.Duration
	now    func() time.Time
}

// New creates a Manager using the provided HMAC secret and TTL.
func New(secret []byte, ttl time.Duration) (*Manager, error) {
	if len(secret) < 16 {
		return nil, errors.New("token: secret must be at least 16 bytes")
	}
	if ttl <= 0 {
		return nil, errors.New("token: ttl must be positive")
	}
	return &Manager{secret: secret, ttl: ttl, now: time.Now}, nil
}

// Issue generates a signed token for the given subject.
func (m *Manager) Issue(subject string) (Token, error) {
	nonce := make([]byte, 8)
	if _, err := rand.Read(nonce); err != nil {
		return Token{}, fmt.Errorf("token: generate nonce: %w", err)
	}
	expiry := m.now().Add(m.ttl).Unix()
	payload := fmt.Sprintf("%s|%d|%s", subject, expiry, hex.EncodeToString(nonce))
	sig := m.sign(payload)
	value := base64.RawURLEncoding.EncodeToString([]byte(payload + "|" + sig))
	return Token{Value: value, ExpiresAt: time.Unix(expiry, 0)}, nil
}

// Validate checks the signature and expiry of a previously issued token.
// It returns the subject on success.
func (m *Manager) Validate(value string) (string, error) {
	raw, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return "", ErrInvalid
	}
	var subject string
	var expiry int64
	var nonce, sig string
	_, err = fmt.Sscanf(string(raw), "%s", new(string))
	// Parse manually to handle subjects with no special chars
	parts := splitLast(string(raw), "|")
	if len(parts) != 2 {
		return "", ErrInvalid
	}
	payload, gotSig := parts[0], parts[1]
	if m.sign(payload) != gotSig {
		return "", ErrInvalid
	}
	_, err = fmt.Sscanf(payload, "%s", new(string))
	n, _ := fmt.Sscanf(payload, "%[^|]|%d|%s", &subject, &expiry, &nonce)
	if n < 2 {
		return "", ErrInvalid
	}
	if m.now().Unix() > expiry {
		return "", ErrExpired
	}
	return subject, nil
}

func (m *Manager) sign(payload string) string {
	h := hmac.New(sha256.New, m.secret)
	h.Write([]byte(payload))
	return hex.EncodeToString(h.Sum(nil))
}

func splitLast(s, sep string) []string {
	idx := -1
	for i := len(s) - 1; i >= 0; i-- {
		if i+len(sep) <= len(s) && s[i:i+len(sep)] == sep {
			idx = i
			break
		}
	}
	if idx == -1 {
		return []string{s}
	}
	return []string{s[:idx], s[idx+len(sep):]}
}
