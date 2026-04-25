package debounce_test

import (
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/yourusername/vaultop/internal/debounce"
)

func TestTrigger_CallsFnAfterDelay(t *testing.T) {
	d := debounce.New(20 * time.Millisecond)

	var mu sync.Mutex
	called := []string{}

	d.Trigger("db/password", func(key string) {
		mu.Lock()
		called = append(called, key)
		mu.Unlock()
	})

	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(called) != 1 || called[0] != "db/password" {
		t.Fatalf("expected fn called once with db/password, got %v", called)
	}
}

func TestTrigger_ResetsTimerOnRepeat(t *testing.T) {
	d := debounce.New(40 * time.Millisecond)

	count := 0
	var mu sync.Mutex

	for i := 0; i < 5; i++ {
		d.Trigger("api/key", func(key string) {
			mu.Lock()
			count++
			mu.Unlock()
		})
		time.Sleep(10 * time.Millisecond)
	}

	time.Sleep(80 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if count != 1 {
		t.Fatalf("expected fn called exactly once, got %d", count)
	}
}

func TestCancel_StopsPendingCall(t *testing.T) {
	d := debounce.New(40 * time.Millisecond)

	called := false
	d.Trigger("secret/x", func(key string) { called = true })
	d.Cancel("secret/x")

	time.Sleep(60 * time.Millisecond)

	if called {
		t.Fatal("expected fn not to be called after Cancel")
	}
}

func TestCancel_NoopWhenNoPending(t *testing.T) {
	d := debounce.New(20 * time.Millisecond)
	// should not panic
	d.Cancel("nonexistent")
}

func TestPending_ReturnsActiveKeys(t *testing.T) {
	d := debounce.New(200 * time.Millisecond)

	d.Trigger("key/a", func(string) {})
	d.Trigger("key/b", func(string) {})

	pending := d.Pending()
	sort.Strings(pending)

	if len(pending) != 2 || pending[0] != "key/a" || pending[1] != "key/b" {
		t.Fatalf("unexpected pending keys: %v", pending)
	}

	d.Cancel("key/a")
	d.Cancel("key/b")
}

func TestPending_EmptyAfterFired(t *testing.T) {
	d := debounce.New(20 * time.Millisecond)
	d.Trigger("k", func(string) {})
	time.Sleep(50 * time.Millisecond)

	if len(d.Pending()) != 0 {
		t.Fatal("expected no pending keys after fn has fired")
	}
}
