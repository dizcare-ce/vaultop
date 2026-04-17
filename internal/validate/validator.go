// Package validate provides secret value validation against configurable rules.
package validate

import (
	"errors"
	"fmt"
	"regexp"
)

// Rule describes a single validation constraint for a secret value.
type Rule struct {
	MinLength int    `yaml:"min_length"`
	MaxLength int    `yaml:"max_length"`
	Pattern   string `yaml:"pattern"`
}

// Result holds the outcome of validating one secret.
type Result struct {
	Key    string
	Passed bool
	Errors []string
}

// Validate checks value against rule and returns a Result.
func Validate(key, value string, rule Rule) Result {
	var errs []string

	if rule.MinLength > 0 && len(value) < rule.MinLength {
		errs = append(errs, fmt.Sprintf("length %d is below minimum %d", len(value), rule.MinLength))
	}

	if rule.MaxLength > 0 && len(value) > rule.MaxLength {
		errs = append(errs, fmt.Sprintf("length %d exceeds maximum %d", len(value), rule.MaxLength))
	}

	if rule.Pattern != "" {
		re, err := regexp.Compile(rule.Pattern)
		if err != nil {
			errs = append(errs, fmt.Sprintf("invalid pattern %q: %v", rule.Pattern, err))
		} else if !re.MatchString(value) {
			errs = append(errs, fmt.Sprintf("value does not match pattern %q", rule.Pattern))
		}
	}

	return Result{Key: key, Passed: len(errs) == 0, Errors: errs}
}

// ValidateAll validates multiple key/value pairs against a map of rules.
// Keys without a rule are skipped. Returns all results and a combined error
// if any validation failed.
func ValidateAll(secrets map[string]string, rules map[string]Rule) ([]Result, error) {
	var results []Result
	var failed int

	for key, rule := range rules {
		val := secrets[key]
		r := Validate(key, val, rule)
		results = append(results, r)
		if !r.Passed {
			failed++
		}
	}

	if failed > 0 {
		return results, errors.New("one or more secrets failed validation")
	}
	return results, nil
}
