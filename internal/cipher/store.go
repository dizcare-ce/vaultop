package cipher

import (
	"context"
	"fmt"
)

// SecretStore is the interface for reading and writing raw secret values.
type SecretStore interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string) error
}

// EncryptedStore wraps a SecretStore, transparently encrypting values on write
// and decrypting on read.
type EncryptedStore struct {
	inner  SecretStore
	cipher *Cipher
}

// NewEncryptedStore returns an EncryptedStore backed by inner.
func NewEncryptedStore(inner SecretStore, c *Cipher) *EncryptedStore {
	return &EncryptedStore{inner: inner, cipher: c}
}

// Set encrypts value and stores it under key.
func (s *EncryptedStore) Set(ctx context.Context, key, value string) error {
	ct, err := s.cipher.Encrypt([]byte(value))
	if err != nil {
		return fmt.Errorf("cipher store set %q: %w", key, err)
	}
	return s.inner.Set(ctx, key, string(ct))
}

// Get retrieves and decrypts the value stored under key.
func (s *EncryptedStore) Get(ctx context.Context, key string) (string, error) {
	raw, err := s.inner.Get(ctx, key)
	if err != nil {
		return "", err
	}
	pt, err := s.cipher.Decrypt([]byte(raw))
	if err != nil {
		return "", fmt.Errorf("cipher store get %q: %w", key, err)
	}
	return string(pt), nil
}
