// Package history provides a lightweight in-memory store for recording
// and querying past cron schedule executions.
//
// Usage:
//
//	store := history.NewStore()
//	err := store.Record(history.Entry{
//		Label:      "backup",
//		Expression: "0 2 * * *",
//		RanAt:      time.Now(),
//		Duration:   3 * time.Second,
//		Success:    true,
//	})
//
// Entries can be retrieved by label, and success rates computed:
//
//	rate := store.SuccessRate("backup")
//
// For rendering, use RenderHistoryTable and RenderSuccessSummary to
// produce human-readable tabular output.
package history
