// Package fingerprint provides content-based fingerprinting for secrets,
// allowing callers to detect whether a secret value has changed without
// storing the plaintext.
package fingerprint

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"
)

// ErrEmptyValue is returned when an empty secret value is fingerprinted.
var ErrEmptyValue = errors.New("fingerprint: value must not be empty")

// Algorithm selects the hashing strategy used to produce a fingerprint.
type Algorithm string

const (
	AlgorithmSHA256 Algorithm = "sha256"
	AlgorithmPrefix Algorithm = "prefix" // first 4 chars + length, for debugging
)

// IsValid reports whether a is a recognised Algorithm.
func (a Algorithm) IsValid() bool {
	switch a {
	case AlgorithmSHA256, AlgorithmPrefix:
		return true
	}
	return false
}

// Of returns a fingerprint string for value using the given algorithm.
// It returns ErrEmptyValue when value is the empty string.
func Of(value string, alg Algorithm) (string, error) {
	if value == "" {
		return "", ErrEmptyValue
	}
	switch alg {
	case AlgorithmSHA256:
		sum := sha256.Sum256([]byte(value))
		return "sha256:" + hex.EncodeToString(sum[:]), nil
	case AlgorithmPrefix:
		runes := []rune(value)
		prefix := string(runes[:min(4, len(runes))])
		return fmt.Sprintf("prefix:%s…(len=%d)", prefix, len(runes)), nil
	default:
		return "", fmt.Errorf("fingerprint: unsupported algorithm %q", alg)
	}
}

// MapOf returns a map of key → fingerprint for every entry in secrets.
// Keys are sorted for deterministic output. Any error short-circuits and
// is returned alongside the partial map built so far.
func MapOf(secrets map[string]string, alg Algorithm) (map[string]string, error) {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	out := make(map[string]string, len(secrets))
	for _, k := range keys {
		fp, err := Of(secrets[k], alg)
		if err != nil {
			return out, fmt.Errorf("fingerprint: key %q: %w", k, err)
		}
		out[k] = fp
	}
	return out, nil
}

// Equal reports whether two fingerprint strings represent the same content.
// Comparison is case-insensitive to tolerate encoding variations.
func Equal(a, b string) bool {
	return strings.EqualFold(a, b)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
