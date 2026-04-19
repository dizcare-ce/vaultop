// Package rollback provides functionality to restore secrets to a previous
// snapshot state, recording the operation in audit and history logs.
package rollback

import (
	"context"
	"fmt"
	"time"

	"github.com/vaultop/internal/audit"
	"github.com/vaultop/internal/history"
	"github.com/vaultop/internal/provider"
	"github.com/vaultop/internal/snapshot"
)

// Options configures a rollback operation.
type Options struct {
	Provider provider.Provider
	History  *history.History
	Audit    audit.Logger
	DryRun   bool
}

// Result holds the outcome of a rollback operation.
type Result struct {
	Key     string
	OldVal  string
	NewVal  string
	Err     error
}

// Run restores all keys from snap into the provider.
// Each key is recorded in history and audit log.
func Run(ctx context.Context, snap snapshot.Snapshot, opts Options) []Result {
	results := make([]Result, 0, len(snap))

	for key, val := range snap {
		r := Result{Key: key, NewVal: val}

		current, err := opts.Provider.Get(ctx, key)
		if err == nil {
			r.OldVal = current
		}

		if !opts.DryRun {
			if err = opts.Provider.Set(ctx, key, val); err != nil {
				r.Err = fmt.Errorf("rollback set %q: %w", key, err)
				_ = opts.History.Record(key, history.Entry{Key: key, Success: false, Timestamp: time.Now(), Note: r.Err.Error()})
				_ = opts.Audit.Log(audit.Entry{Action: "rollback", Key: key, Success: false, Error: r.Err.Error()})
				results = append(results, r)
				continue
			}
			_ = opts.History.Record(key, history.Entry{Key: key, Success: true, Timestamp: time.Now(), Note: "rollback"})
			_ = opts.Audit.Log(audit.Entry{Action: "rollback", Key: key, Success: true})
		}

		results = append(results, r)
	}

	return results
}

// AnyFailed returns true if any result contains an error.
func AnyFailed(results []Result) bool {
	for _, r := range results {
		if r.Err != nil {
			return true
		}
	}
	return false
}
