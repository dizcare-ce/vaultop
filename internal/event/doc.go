// Package event implements a lightweight in-process publish/subscribe bus used
// to broadcast secret lifecycle events (rotation, deletion, import, expiry)
// to interested components without creating direct dependencies between them.
//
// Usage:
//
//	bus := event.New()
//	bus.Subscribe(event.KindRotated, func(e event.Event) {
//		fmt.Println("rotated:", e.Key)
//	})
//	bus.Publish(event.Event{Kind: event.KindRotated, Key: "db/password"})
package event
