package webhook_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/vaultop/internal/webhook"
)

func TestSend_Success(t *testing.T) {
	var received webhook.Event
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode: %v", err)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	s := webhook.New(srv.URL, 5*time.Second)
	e := webhook.Event{Kind: "rotated", Key: "db/pass", Meta: map[string]string{"env": "prod"}}
	if err := s.Send(context.Background(), e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Kind != "rotated" {
		t.Errorf("kind: got %q, want %q", received.Kind, "rotated")
	}
	if received.Key != "db/pass" {
		t.Errorf("key: got %q, want %q", received.Key, "db/pass")
	}
}

func TestSend_NonSuccessStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	s := webhook.New(srv.URL, 5*time.Second)
	err := s.Send(context.Background(), webhook.Event{Kind: "rotated", Key: "x"})
	if err == nil {
		t.Fatal("expected error for 500 response")
	}
}

func TestSend_TimestampAutoFilled(t *testing.T) {
	var received webhook.Event
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	s := webhook.New(srv.URL, 5*time.Second)
	s.Send(context.Background(), webhook.Event{Kind: "test", Key: "k"})
	if received.Timestamp.IsZero() {
		t.Error("expected timestamp to be set automatically")
	}
}

func TestNoop_Send(t *testing.T) {
	var n webhook.Noop
	if err := n.Send(context.Background(), webhook.Event{Kind: "rotated", Key: "k"}); err != nil {
		t.Fatalf("noop should not error: %v", err)
	}
}
