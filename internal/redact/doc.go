// Package redact provides helpers for masking sensitive secret values
// before they appear in logs, audit entries, notifications, or CLI output.
//
// Three redaction modes are supported:
//
//   - ModeFull    – replaces the entire value with a mask string (default "***")
//   - ModePartial – preserves a configurable suffix so operators can
//                   identify a secret without exposing it fully
//   - ModeHash    – replaces every character with "*", preserving length
//
// Use RedactMap to bulk-redact all values in a map before passing secrets
// to any output-producing subsystem.
package redact
