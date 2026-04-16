// Package export serialises secret maps into portable output formats.
//
// Supported formats:
//
//   - json    — pretty-printed JSON object
//   - env     — KEY=VALUE lines suitable for shell sourcing
//   - dotenv  — alias for env, compatible with .env file conventions
//
// Usage:
//
//	secrets := map[string]string{"DB_PASS": "s3cr3t", "API_KEY": "abc"}
//	if err := export.Write(os.Stdout, secrets, export.FormatEnv); err != nil {
//		log.Fatal(err)
//	}
package export
