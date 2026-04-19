package cipher_test

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/vaultop/internal/cipher"
)

type memStore struct {
	mu   sync.Mutex
	data map[string]string
}

func newMemStore() *memStore { return &memStore{data: map[string]string{}} }

func (m *memStore) Get(_ context.Context, key string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.data[key]
	if !ok {
		return "", errors.New("not found")
	}
	return v, nil
}

func (m *memStore) Set(_ context.Context, key, value string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
	return nil
}

func TestEncryptedStore_RoundTrip(t *testing.T) {
	c, _ := cipher.New(key32(), cipher.AES256GCM)
	store := cipher.NewEncryptedStore(newMemStore(), c)
	ctx := context.Background()

	if err := store.Set(ctx, "db/pass", "hunter2"); err != nil {
		t.Fatalf("set: %v", err)
	}
	got, err := store.Get(ctx, "db/pass")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got != "hunter2" {
		t.Fatalf("got %q, want %q", got, "hunter2")
	}
}

func TestEncryptedStore_StoredValueIsNotPlaintext(t *testing.T) {
	c, _ := cipher.New(key32(), cipher.AES256GCM)
	mem := newMemStore()
	store := cipher.NewEncryptedStore(mem, c)
	ctx := context.Background()

	_ = store.Set(ctx, "key", "plaintext")
	raw, _ := mem.Get(ctx, "key")
	if raw == "plaintext" {
		t.Fatal("expected stored value to be encrypted")
	}
}

func TestEncryptedStore_GetMissingKey(t *testing.T) {
	c, _ := cipher.New(key32(), cipher.AES256GCM)
	store := cipher.NewEncryptedStore(newMemStore(), c)
	_, err := store.Get(context.Background(), "missing")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}
