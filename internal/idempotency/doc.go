// Package idempotency provides a time-windowed deduplication store for
// secret operations in vaultop.
//
// When a caller retries a rotation, import, or write operation it can supply
// a stable idempotency key (e.g. a client-generated UUID). The Store records
// the outcome of the first attempt and returns the cached Result for any
// subsequent attempt that arrives within the configured window, preventing
// duplicate side-effects such as double-rotation or double-import.
//
// Usage:
//
//	store := idempotency.New(10 * time.Minute)
//
//	if r, ok := store.Check(key); ok {
//	    return r.Value, r.Err // replay cached outcome
//	}
//
//	value, err := doExpensiveOperation()
//	store.Record(key, idempotency.Result{Value: value, Err: err})
package idempotency
