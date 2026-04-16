package audit

import "io"

// Noop returns a Logger that discards all output.
// Useful in tests or when auditing is disabled.
func Noop() *Logger {
	return NewLogger(io.Discard)
}
