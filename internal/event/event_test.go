package event_test

import (
	"testing"

	"github.com/vaultop/internal/event"
)

func TestPublish_DeliversToSubscriber(t *testing.T) {
	bus := event.New()
	var got event.Event
	bus.Subscribe(event.KindRotated, func(e event.Event) { got = e })

	bus.Publish(event.Event{Kind: event.KindRotated, Key: "db/pass", Message: "rotated"})

	if got.Key != "db/pass" {
		t.Fatalf("expected key db/pass, got %q", got.Key)
	}
}

func TestPublish_MultipleSubscribers(t *testing.T) {
	bus := event.New()
	count := 0
	inc := func(event.Event) { count++ }
	bus.Subscribe(event.KindDeleted, inc)
	bus.Subscribe(event.KindDeleted, inc)

	bus.Publish(event.Event{Kind: event.KindDeleted, Key: "k"})

	if count != 2 {
		t.Fatalf("expected 2 calls, got %d", count)
	}
}

func TestPublish_WrongKind_NotDelivered(t *testing.T) {
	bus := event.New()
	called := false
	bus.Subscribe(event.KindImported, func(event.Event) { called = true })

	bus.Publish(event.Event{Kind: event.KindRotated, Key: "k"})

	if called {
		t.Fatal("handler should not have been called for a different kind")
	}
}

func TestReset_ClearsHandlers(t *testing.T) {
	bus := event.New()
	called := false
	bus.Subscribe(event.KindExpired, func(event.Event) { called = true })
	bus.Reset()

	bus.Publish(event.Event{Kind: event.KindExpired, Key: "k"})

	if called {
		t.Fatal("handler should have been removed after Reset")
	}
}

func TestPublish_NoSubscribers_NoPanic(t *testing.T) {
	bus := event.New()
	// should not panic
	bus.Publish(event.Event{Kind: event.KindRotated, Key: "k"})
}
