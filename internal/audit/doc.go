// Package audit implements structured audit logging for vaultop.
//
// Every secret rotation — whether live, dry-run, or failed — should produce
// an audit entry so operators can trace what changed, when, and why.
//
// Entries are written as newline-delimited JSON (NDJSON) to any io.Writer,
// making them easy to ship to log aggregators or store on disk.
//
// Basic usage:
//
//	l := audit.NewLogger(os.Stdout)
//	l.Log(audit.EventRotated, "aws", "prod/db/password", "")
//
// For file-backed logging:
//
//	l, err := audit.NewFileLogger("/var/log/vaultop/audit.log")
package audit
