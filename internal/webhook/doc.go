// Package webhook provides HTTP webhook delivery for vaultop events.
//
// Use New to create a Sender pointed at a remote endpoint. Events are
// serialised as JSON and delivered via HTTP POST. A Noop implementation
// is available for testing or when webhook delivery is disabled.
//
// Example:
//
//	s := webhook.New("https://example.com/hook", 5*time.Second)
//	err := s.Send(ctx, webhook.Event{Kind: "rotated", Key: "db/password"})
package webhook
