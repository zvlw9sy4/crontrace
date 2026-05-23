// Package visualizer provides terminal-based rendering of cron schedule data.
package visualizer

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// ScheduleRow represents a single row in the schedule table.
type ScheduleRow struct {
	Expression string
	Description string
	NextRuns    []time.Time
}

// ConflictRow represents a detected conflict between two schedules.
type ConflictRow struct {
	ExpressionA string
	ExpressionB string
	At          time.Time
}

// RenderScheduleTable writes a formatted ASCII table of schedule rows to w.
func RenderScheduleTable(w io.Writer, rows []ScheduleRow) {
	const colExpr = 30
	const colDesc = 25
	const colNext = 22

	header := fmt.Sprintf("%-*s %-*s %s", colExpr, "EXPRESSION", colDesc, "DESCRIPTION", "NEXT RUNS")
	sep := strings.Repeat("-", len(header)+10)

	fmt.Fprintln(w, sep)
	fmt.Fprintln(w, header)
	fmt.Fprintln(w, sep)

	for _, row := range rows {
		nextStrs := make([]string, len(row.NextRuns))
		for i, t := range row.NextRuns {
			nextStrs[i] = t.Format("2006-01-02 15:04:05")
		}
		nextCol := strings.Join(nextStrs, "  |  ")
		fmt.Fprintf(w, "%-*s %-*s %s\n", colExpr, row.Expression, colDesc, row.Description, nextCol)
	}

	fmt.Fprintln(w, sep)
}

// RenderConflictTable writes a formatted ASCII table of conflict rows to w.
func RenderConflictTable(w io.Writer, rows []ConflictRow) {
	if len(rows) == 0 {
		fmt.Fprintln(w, "No conflicts detected.")
		return
	}

	const colA = 28
	const colB = 28

	header := fmt.Sprintf("%-*s %-*s %s", colA, "EXPRESSION A", colB, "EXPRESSION B", "CONFLICT AT")
	sep := strings.Repeat("-", len(header)+4)

	fmt.Fprintln(w, "⚠ Conflicts Detected:")
	fmt.Fprintln(w, sep)
	fmt.Fprintln(w, header)
	fmt.Fprintln(w, sep)

	for _, row := range rows {
		fmt.Fprintf(w, "%-*s %-*s %s\n", colA, row.ExpressionA, colB, row.ExpressionB, row.At.Format("2006-01-02 15:04:05"))
	}

	fmt.Fprintln(w, sep)
}
