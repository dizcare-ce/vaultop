package pipeline_test

import (
	"context"
	"errors"
	"testing"

	"github.com/vaultop/internal/pipeline"
)

func step(name string, err error) pipeline.Step {
	return pipeline.Step{
		Name: name,
		Fn: func(_ context.Context, s *pipeline.State) error {
			if err == nil {
				s.Values[name] = "done"
			}
			return err
		},
	}
}

func TestRun_AllSucceed(t *testing.T) {
	p := pipeline.New(step("a", nil), step("b", nil))
	state := pipeline.NewState()
	results, err := p.Run(context.Background(), state)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if state.Values["a"] != "done" || state.Values["b"] != "done" {
		t.Error("state not updated by steps")
	}
}

func TestRun_StopsOnError(t *testing.T) {
	errBoom := errors.New("boom")
	p := pipeline.New(step("a", nil), step("b", errBoom), step("c", nil))
	results, err := p.Run(context.Background(), pipeline.NewState())
	if err == nil {
		t.Fatal("expected error")
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestRun_ContinueOnError(t *testing.T) {
	errBoom := errors.New("boom")
	p := pipeline.New(step("a", errBoom), step("b", nil))
	p.ContinueOnError = true
	results, err := p.Run(context.Background(), pipeline.NewState())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Err == nil {
		t.Error("expected step a to record error")
	}
}

func TestRun_EmptyPipeline(t *testing.T) {
	p := pipeline.New()
	results, err := p.Run(context.Background(), pipeline.NewState())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Error("expected no results")
	}
}
