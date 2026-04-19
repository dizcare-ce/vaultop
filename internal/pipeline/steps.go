package pipeline

import (
	"context"
	"fmt"

	"github.com/vaultop/internal/provider"
	"github.com/vaultop/internal/validate"
)

// FetchStep returns a Step that loads secrets from a provider into state.
func FetchStep(p provider.Provider, keys []string) Step {
	return Step{
		Name: "fetch",
		Fn: func(ctx context.Context, s *State) error {
			for _, k := range keys {
				v, err := p.Get(ctx, k)
				if err != nil {
					return fmt.Errorf("fetch %q: %w", k, err)
				}
				s.Values[k] = v
			}
			return nil
		},
	}
}

// ValidateStep returns a Step that validates fetched secrets against rules.
func ValidateStep(rules []validate.Rule) Step {
	return Step{
		Name: "validate",
		Fn: func(_ context.Context, s *State) error {
			for k, v := range s.Values {
				if err := validate.Validate(k, v, rules); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

// WriteStep returns a Step that persists state values back to a provider.
func WriteStep(p provider.Provider) Step {
	return Step{
		Name: "write",
		Fn: func(ctx context.Context, s *State) error {
			for k, v := range s.Values {
				if err := p.Set(ctx, k, v); err != nil {
					return fmt.Errorf("write %q: %w", k, err)
				}
			}
			return nil
		},
	}
}
