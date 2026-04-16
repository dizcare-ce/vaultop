package rotation

import (
	"context"
	"fmt"
	"time"

	"github.com/vaultop/internal/provider"
)

// Result holds the outcome of a single secret rotation.
type Result struct {
	SecretID  string
	Provider  string
	RotatedAt time.Time
	Err       error
}

// Options configures rotation behaviour.
type Options struct {
	DryRun    bool
	Generator ValueGenerator
}

// Rotator rotates secrets using the given provider.
type Rotator struct {
	p   provider.Provider
	opt Options
}

// New creates a Rotator for the supplied provider and options.
func New(p provider.Provider, opt Options) *Rotator {
	if opt.Generator == nil {
		opt.Generator = DefaultGenerator
	}
	return &Rotator{p: p, opt: opt}
}

// Rotate replaces the value of each secretID with a freshly generated value.
// It returns one Result per secret.
func (r *Rotator) Rotate(ctx context.Context, secretIDs []string) []Result {
	results := make([]Result, 0, len(secretIDs))
	for _, id := range secretIDs {
		res := Result{SecretID: id, Provider: string(r.p.Type()), RotatedAt: time.Now()}
		newVal, err := r.opt.Generator(id)
		if err != nil {
			res.Err = fmt.Errorf("generate: %w", err)
			results = append(results, res)
			continue
		}
		if !r.opt.DryRun {
			if err := r.p.Set(ctx, id, newVal); err != nil {
				res.Err = fmt.Errorf("set: %w", err)
			}
		}
		results = append(results, res)
	}
	return results
}
