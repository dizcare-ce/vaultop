// Package compress provides optional compression for secret values before
// storing them in a provider backend.
package compress

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
)

// Algorithm identifies a supported compression algorithm.
type Algorithm string

const (
	Gzip Algorithm = "gzip"
	None Algorithm = "none"
)

// IsValid reports whether the algorithm is supported.
func (a Algorithm) IsValid() bool {
	switch a {
	case Gzip, None:
		return true
	}
	return false
}

// Compress compresses src using the given algorithm and returns a
// base64-encoded string suitable for storage.
func Compress(alg Algorithm, src []byte) (string, error) {
	if alg == None {
		return base64.StdEncoding.EncodeToString(src), nil
	}
	if alg != Gzip {
		return "", fmt.Errorf("compress: unsupported algorithm %q", alg)
	}
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	if _, err := w.Write(src); err != nil {
		return "", fmt.Errorf("compress: gzip write: %w", err)
	}
	if err := w.Close(); err != nil {
		return "", fmt.Errorf("compress: gzip close: %w", err)
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// Decompress reverses Compress, returning the original plaintext bytes.
func Decompress(alg Algorithm, src string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return nil, fmt.Errorf("compress: base64 decode: %w", err)
	}
	if alg == None {
		return decoded, nil
	}
	if alg != Gzip {
		return nil, fmt.Errorf("compress: unsupported algorithm %q", alg)
	}
	r, err := gzip.NewReader(bytes.NewReader(decoded))
	if err != nil {
		return nil, fmt.Errorf("compress: gzip reader: %w", err)
	}
	defer r.Close()
	out, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("compress: gzip read: %w", err)
	}
	return out, nil
}
