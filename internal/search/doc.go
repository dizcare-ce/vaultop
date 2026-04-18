// Package search implements secret lookup and filtering across a provider.
//
// It supports filtering by key prefix, key substring, and value substring,
// allowing operators to quickly locate secrets without knowing exact names.
//
// Example:
//
//	res, err := search.Find(ctx, provider, search.Options{
//		Prefix:   "app/",
//		Contains: "db",
//	})
package search
