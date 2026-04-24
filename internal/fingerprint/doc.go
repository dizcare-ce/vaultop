// Package fingerprint provides lightweight, algorithm-agnostic fingerprinting
// of secret values for vaultop.
//
// A fingerprint is a short, opaque string derived from a secret's content.
// It lets callers determine whether a secret has changed since it was last
// observed — without retaining or comparing the plaintext.
//
// Supported algorithms:
//
//	"sha256"  — SHA-256 hex digest, prefixed with "sha256:"
//	"prefix"  — first four runes plus the total rune-length; intended for
//	            human-readable debugging only, not for security comparisons.
//
// Typical usage:
//
//	fp, err := fingerprint.Of(secretValue, fingerprint.AlgorithmSHA256)
//	if err != nil { ... }
//	if !fingerprint.Equal(fp, storedFingerprint) {
//	    // secret has rotated
//	}
package fingerprint
