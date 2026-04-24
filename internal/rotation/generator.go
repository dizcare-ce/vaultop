package rotation

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// ValueGenerator produces a new secret value for the given secretID.
type ValueGenerator func(secretID string) (string, error)

// DefaultGenerator returns a 32-byte cryptographically random base64 string.
var DefaultGenerator ValueGenerator = func(_ string) (string, error) {
	return randomBase64(32)
}

// FixedGenerator always returns the same value; useful for testing.
func FixedGenerator(value string) ValueGenerator {
	return func(_ string) (string, error) {
		return value, nil
	}
}

// ErrorGenerator always returns an error; useful for testing.
func ErrorGenerator(msg string) ValueGenerator {
	return func(_ string) (string, error) {
		return "", fmt.Errorf("%s", msg)
	}
}

// RandomBytesGenerator returns a ValueGenerator that produces a cryptographically
// random base64 string of the given byte length. It panics if n <= 0.
func RandomBytesGenerator(n int) ValueGenerator {
	if n <= 0 {
		panic(fmt.Sprintf("rotation: RandomBytesGenerator requires n > 0, got %d", n))
	}
	return func(_ string) (string, error) {
		return randomBase64(n)
	}
}

func randomBase64(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("rand.Read: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
