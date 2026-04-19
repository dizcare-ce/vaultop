package compress_test

import (
	"strings"
	"testing"

	"github.com/vaultop/internal/compress"
)

func TestAlgorithm_IsValid(t *testing.T) {
	if !compress.Gzip.IsValid() {
		t.Error("expected Gzip to be valid")
	}
	if !compress.None.IsValid() {
		t.Error("expected None to be valid")
	}
	if compress.Algorithm("lz4").IsValid() {
		t.Error("expected lz4 to be invalid")
	}
}

func TestCompressDecompress_Gzip_RoundTrip(t *testing.T) {
	original := []byte("super-secret-value-that-is-long-enough-to-compress-well-aaaaaaa")
	encoded, err := compress.Compress(compress.Gzip, original)
	if err != nil {
		t.Fatalf("Compress: %v", err)
	}
	got, err := compress.Decompress(compress.Gzip, encoded)
	if err != nil {
		t.Fatalf("Decompress: %v", err)
	}
	if string(got) != string(original) {
		t.Errorf("round-trip mismatch: got %q want %q", got, original)
	}
}

func TestCompressDecompress_None_RoundTrip(t *testing.T) {
	original := []byte("plain-value")
	encoded, err := compress.Compress(compress.None, original)
	if err != nil {
		t.Fatalf("Compress: %v", err)
	}
	got, err := compress.Decompress(compress.None, encoded)
	if err != nil {
		t.Fatalf("Decompress: %v", err)
	}
	if string(got) != string(original) {
		t.Errorf("round-trip mismatch: got %q want %q", got, original)
	}
}

func TestCompress_UnsupportedAlgorithm(t *testing.T) {
	_, err := compress.Compress("lz4", []byte("data"))
	if err == nil || !strings.Contains(err.Error(), "unsupported") {
		t.Errorf("expected unsupported error, got %v", err)
	}
}

func TestDecompress_UnsupportedAlgorithm(t *testing.T) {
	encoded, _ := compress.Compress(compress.None, []byte("data"))
	_, err := compress.Decompress("lz4", encoded)
	if err == nil || !strings.Contains(err.Error(), "unsupported") {
		t.Errorf("expected unsupported error, got %v", err)
	}
}

func TestDecompress_InvalidBase64(t *testing.T) {
	_, err := compress.Decompress(compress.Gzip, "!!!not-base64!!!")
	if err == nil {
		t.Error("expected base64 decode error")
	}
}
