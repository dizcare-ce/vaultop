// Package ttl provides secret expiry tracking and expiration checks.
package ttl

import (
	"errors"
	"time"
)

// ErrExpired is returned when a secret has passed its expiry time.
var ErrExpired = errors.New("secret has expired")

// Entry holds expiry metadata for a single secret key.
type Entry struct {
	Key       string    `json:"key"`
	ExpiresAt time.Time `json:"expires_at"`
}

// IsExpired reports whether the entry has expired relative to now.
func (e Entry) IsExpired(now time.Time) bool {
	return now.After(e.ExpiresAt)
}

// TTLMap manages expiry entries keyed by secret name.
type TTLMap struct {
	entries map[string]Entry
}

// New returns an empty TTLMap.
func New() *TTLMap {
	return &TTLMap{entries: make(map[string]Entry)}
}

// Set registers an expiry duration for the given key, starting from now.
func (m *TTLMap) Set(key string, ttl time.Duration, now time.Time) {
	m.entries[key] = Entry{
		Key:       key,
		ExpiresAt: now.Add(ttl),
	}
}

// Get returns the Entry for key and whether it was found.
func (m *TTLMap) Get(key string) (Entry, bool) {
	e, ok := m.entries[key]
	return e, ok
}

// Delete removes the expiry entry for key.
func (m *TTLMap) Delete(key string) {
	delete(m.entries, key)
}

// Expired returns all keys that have expired relative to now.
func (m *TTLMap) Expired(now time.Time) []string {
	var out []string
	for k, e := range m.entries {
		if e.IsExpired(now) {
			out = append(out, k)
		}
	}
	return out
}

// Check returns ErrExpired if the key exists and has expired.
func (m *TTLMap) Check(key string, now time.Time) error {
	e, ok := m.entries[key]
	if !ok {
		return nil
	}
	if e.IsExpired(now) {
		return ErrExpired
	}
	return nil
}
