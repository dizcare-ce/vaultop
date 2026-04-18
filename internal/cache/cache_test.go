package cache_test

import (
	"testing"
	"time"

	"github.com/vaultop/internal/cache"
)

func TestSet_And_Get(t *testing.T) {
	c := cache.New(5 * time.Second)
	c.Set("foo", "bar")
	v, ok := c.Get("foo")
	if !ok {
		t.Fatal("expected hit")
	}
	if v != "bar" {
		t.Fatalf("got %q, want %q", v, "bar")
	}
}

func TestGet_Miss(t *testing.T) {
	c := cache.New(5 * time.Second)
	_, ok := c.Get("missing")
	if ok {
		t.Fatal("expected miss for unknown key")
	}
}

func TestGet_Expired(t *testing.T) {
	c := cache.New(10 * time.Millisecond)
	c.Set("key", "val")
	time.Sleep(20 * time.Millisecond)
	_, ok := c.Get("key")
	if ok {
		t.Fatal("expected expired entry to be a miss")
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	c := cache.New(5 * time.Second)
	c.Set("k", "v")
	c.Delete("k")
	_, ok := c.Get("k")
	if ok {
		t.Fatal("expected miss after delete")
	}
}

func TestFlush_ClearsAll(t *testing.T) {
	c := cache.New(5 * time.Second)
	c.Set("a", "1")
	c.Set("b", "2")
	c.Flush()
	if c.Len() != 0 {
		t.Fatalf("expected 0 entries after flush, got %d", c.Len())
	}
}

func TestLen_CountsEntries(t *testing.T) {
	c := cache.New(5 * time.Second)
	c.Set("x", "1")
	c.Set("y", "2")
	if c.Len() != 2 {
		t.Fatalf("expected 2, got %d", c.Len())
	}
}
