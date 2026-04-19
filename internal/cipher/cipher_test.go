package cipher_test

import (
	"bytes"
	"testing"

	"github.com/vaultop/internal/cipher"
)

func key32() []byte {
	return bytes.Repeat([]byte("k"), 32)
}

func TestNew_ValidKey(t *testing.T) {
	_, err := cipher.New(key32(), cipher.AES256GCM)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_InvalidAlgorithm(t *testing.T) {
	_, err := cipher.New(key32(), cipher.Algorithm("rot13"))
	if err == nil {
		t.Fatal("expected error for unsupported algorithm")
	}
}

func TestNew_ShortKey(t *testing.T) {
	_, err := cipher.New([]byte("short"), cipher.AES256GCM)
	if err == nil {
		t.Fatal("expected error for short key")
	}
}

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	c, _ := cipher.New(key32(), cipher.AES256GCM)
	plaintext := []byte("super-secret-value")
	ct, err := c.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	got, err := c.Decrypt(ct)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	if !bytes.Equal(got, plaintext) {
		t.Fatalf("got %q, want %q", got, plaintext)
	}
}

func TestEncrypt_ProducesUniqueNonces(t *testing.T) {
	c, _ := cipher.New(key32(), cipher.AES256GCM)
	a, _ := c.Encrypt([]byte("value"))
	b, _ := c.Encrypt([]byte("value"))
	if bytes.Equal(a, b) {
		t.Fatal("expected unique ciphertexts")
	}
}

func TestDecrypt_TamperedCiphertext(t *testing.T) {
	c, _ := cipher.New(key32(), cipher.AES256GCM)
	ct, _ := c.Encrypt([]byte("value"))
	ct[len(ct)-1] ^= 0xff
	if _, err := c.Decrypt(ct); err == nil {
		t.Fatal("expected error for tampered ciphertext")
	}
}

func TestDecrypt_TooShort(t *testing.T) {
	c, _ := cipher.New(key32(), cipher.AES256GCM)
	if _, err := c.Decrypt([]byte{0x01}); err == nil {
		t.Fatal("expected error for short ciphertext")
	}
}
