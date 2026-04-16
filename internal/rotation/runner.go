package rotation

import (
	"context"
	"fmt"
	"time"

	"github.com/yourusername/vaultop/internal/audit"
	"github.com/yourusername/vaultop/internal/history"
	"github.com/yourusername/vaultop/internal/notification"
	"github.com/yourusername/vaultop/internal/provider"
)

// RunOptions holds dependencies and configuration for a rotation run.
type RunOptions struct {
	Provider     provider.Provider
	Audit        audit.Logger
	History      *history.History
	Notifier     notification.Notifier
	Generator    Generator
	DryRun       bool
}

// RunResult summarises the outcome of rotating a set of policies.
type RunResult struct {
	Rotated []string
	Skipped []string
	Failed  map[string]error
}

// Run executes rotation for each policy, recording history and emitting
// audit log entries and notifications for every key.
func Run(ctx context.Context, policies []Policy, opts RunOptions) (RunResult, error) {
	result := RunResult{Failed: make(map[string]error)}

	rotator := New(opts.Provider, opts.Generator)

	for _, p := range policies {
		if err := p.Validate(); err != nil {
			result.Failed[p.Key] = fmt.Errorf("invalid policy: %w", err)
			continue
		}

		start := time.Now()
		err := rotator.Rotate(ctx, p, opts.DryRun)
		duration := time.Since(start)

		event := notification.Event{
			Key:      p.Key,
			DryRun:   opts.DryRun,
			Duration: duration,
			Err:      err,
		}

		_ = opts.Notifier.Notify(ctx, event)

		if err != nil {
			result.Failed[p.Key] = err
			_ = opts.Audit.Log(audit.Entry{Key: p.Key, Action: "rotate", Success: false, Error: err.Error()})
			if !opts.DryRun {
				_ = opts.History.Record(p.Key, history.StatusFailure, err.Error())
			}
			continue
		}

		result.Rotated = append(result.Rotated, p.Key)
		_ = opts.Audit.Log(audit.Entry{Key: p.Key, Action: "rotate", Success: true})
		if !opts.DryRun {
			_ = opts.History.Record(p.Key, history.StatusSuccess, "")
		}
	}

	return result, nil
}
