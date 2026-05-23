// Package snapshot provides save/load functionality for cron schedule snapshots.
//
// A snapshot captures the current state of one or more cron schedule entries —
// including their expressions, predicted next-run times, and the timestamp at
// which the snapshot was taken. Snapshots are persisted as JSON files and can
// be reloaded for diffing, auditing, or historical comparison.
//
// Typical usage:
//
//	entries := []snapshot.Entry{
//		{Label: "backup", Expression: "0 2 * * *", NextRuns: runs, SavedAt: time.Now()},
//	}
//	s := snapshot.New(entries)
//	if err := snapshot.Save(s, "./backup.snap.json"); err != nil {
//		log.Fatal(err)
//	}
//
//	loaded, err := snapshot.Load("./backup.snap.json")
//	if err != nil {
//		log.Fatal(err)
//	}
//	snapshot.RenderSnapshotTable(os.Stdout, loaded)
package snapshot
