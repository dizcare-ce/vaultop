// Package pipeline provides a sequential multi-step secret operation executor.
package pipeline

import (
	"context"
	"fmt"
	"time"
)

// StepFunc is a single pipeline step.
type StepFunc func(ctx context.Context, state *State) error

// Step wraps a named StepFunc.
type Step struct {
	Name string
	Fn   StepFunc
}

// State carries shared data between pipeline steps.
type State struct {
	Values map[string]string
	Meta   map[string]any
}

func NewState() *State {
	return &State{
		Values: make(map[string]string),
		Meta:   make(map[string]any),
	}
}

// Result holds the outcome of a single step.
type Result struct {
	Step     string
	Duration time.Duration
	Err      error
}

// Pipeline executes steps in order, stopping on first error unless ContinueOnError is set.
type Pipeline struct {
	steps           []Step
	ContinueOnError bool
}

func New(steps ...Step) *Pipeline {
	return &Pipeline{steps: steps}
}

// Run executes all steps and returns per-step results.
func (p *Pipeline) Run(ctx context.Context, state *State) ([]Result, error) {
	results := make([]Result, 0, len(p.steps))
	for _, s := range p.steps {
		start := time.Now()
		err := s.Fn(ctx, state)
		results = append(results, Result{Step: s.Name, Duration: time.Since(start), Err: err})
		if err != nil && !p.ContinueOnError {
			return results, fmt.Errorf("pipeline step %q failed: %w", s.Name, err)
		}
	}
	return results, nil
}
