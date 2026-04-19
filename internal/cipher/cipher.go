// Package cipher provides symmetric encryption and decryption for secret values.
package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

// Algorithm represents a supported encryption algorithm.
type Algorithm string

const (
	AES256GCM Algorithm = "aes256gcm"
)

// IsValid reports whether the algorithm is supported.
func (a Algorithm) IsValid() bool {
	return a == AES256GCM
}

// Cipher encrypts and decrypts values using a symmetric key.
type Cipher struct {
	key []byte
	alg Algorithm
}

// New returns a Cipher for the given 32-byte key and algorithm.
func New(key []byte, alg Algorithm) (*Cipher, error) {
	if !alg.IsValid() {
		return nil, errors.New("cipher: unsupported algorithm: " + string(alg))
	}
	if len(key) != 32 {
		return nil, errors.New("cipher: key must be 32 bytes")
	}
	return &Cipher{key: key, alg: alg}, nil
}

// Encrypt encrypts plaintext and returns ciphertext with a prepended nonce.
func (c *Cipher) Encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// Decrypt decrypts ciphertext produced by Encrypt.
func (c *Cipher) Decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	ns := gcm.NonceSize()
	if len(ciphertext) < ns {
		return nil, errors.New("cipher: ciphertext too short")
	}
	return gcm.Open(nil, ciphertext[:ns], ciphertext[ns:], nil)
}
