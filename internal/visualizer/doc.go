// Package visualizer renders cron schedule data as human-readable terminal
// output. It provides two main capabilities:
//
//  1. Table rendering – ASCII tables showing expressions, descriptions, and
//     upcoming run times via RenderScheduleTable and RenderConflictTable.
//
//  2. Description generation – Describe converts a standard five-field cron
//     expression into a plain-English label (e.g. "Every hour", "Daily at
//     midnight"). Unknown patterns fall back to "Custom schedule".
//
// Typical usage:
//
//	rows := []visualizer.ScheduleRow{
//		{
//			Expression:  "0 * * * *",
//			Description: visualizer.Describe("0 * * * *"),
//			NextRuns:    runs,
//		},
//	}
//	visualizer.RenderScheduleTable(os.Stdout, rows)
//
// The package writes to any io.Writer, making it straightforward to capture
// output in tests or redirect it to a file.
package visualizer
