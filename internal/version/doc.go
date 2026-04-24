// Package version tracks the history of secret values as numbered versions.
//
// Each call to Record appends a new Entry for the given key, incrementing
// the version counter. Callers can retrieve a specific version with Get,
// the most recent value with Latest, or the full history with List.
//
// The Store is safe for concurrent use.
//
// Example usage:
//
//	store := version.New()
//	store.Record("api/key", "old-value")
//	store.Record("api/key", "new-value")
//
//	latest, err := store.Latest("api/key")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(latest.Version, latest.Value)
package version
