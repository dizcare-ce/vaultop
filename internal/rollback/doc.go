// Package rollback restores provider secrets to a previously captured
// snapshot state.
//
// Typical usage:
//
//	snap, _ := snapshot.Load("backup.json")
//	results := rollback.Run(ctx, snap, rollback.Options{
//		Provider: p,
//		History:  h,
//		Audit:    logger,
//	})
//	if rollback.AnyFailed(results) {
//		// handle errors
//	}
//
// Dry-run mode previews which keys would be restored without writing
// any changes to the provider or recording history entries.
package rollback
