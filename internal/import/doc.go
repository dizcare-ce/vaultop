// Package importer handles importing secrets from external files into
// a vaultop-managed provider.
//
// Supported formats:
//
//   - json: a flat JSON object mapping string keys to string values.
//   - env:  a KEY=VALUE file, one entry per line; lines starting with
//     '#' and blank lines are ignored.
//
// Example (JSON):
//
//	imported, err := importer.Import(r, importer.Options{
//		Provider: p,
//		Format:   importer.FormatJSON,
//	})
//
// Example (env file):
//
//	imported, err := importer.ImportFile(".env", importer.Options{
//		Provider: p,
//		Format:   importer.FormatEnv,
//		DryRun:   true,
//	})
package importer
