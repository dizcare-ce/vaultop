// Package lock provides file-based mutual exclusion for vaultop operations.
//
// A Lock prevents concurrent rotation runs from interfering with each other
// by creating an exclusive lock file that records the owning process ID.
//
// Typical usage:
//
//	l := lock.New("/tmp/vaultop.lock")
//	if err := l.Acquire(); err != nil {
//		// another process is running
//		return err
//	}
//	defer l.Release()
//
// The lock file is automatically removed on Release. If the process crashes
// the file must be cleaned up manually or via a startup check.
package lock
