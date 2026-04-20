// Package resolver provides secret key resolution with alias expansion
// and ordered fallback support.
//
// A Resolver wraps a Provider and applies a two-layer lookup strategy:
//
//  1. Alias expansion — a logical name (e.g. "db_pass") is mapped to its
//     canonical provider key before the first lookup attempt.
//
//  2. Fallback chain — if the primary key is absent, the resolver tries each
//     fallback key in order, returning the first successful value.
//
// Example:
//
//	cfg := resolver.Config{
//		Aliases:   map[string]string{"db_pass": "prod/db/password"},
//		Fallbacks: map[string][]string{"db_pass": {"staging/db/password"}},
//	}
//	r := resolver.New(myProvider, cfg)
//	val, err := r.Resolve("db_pass")
package resolver
