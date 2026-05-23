package history

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

// RenderHistoryTable writes a formatted table of history entries to w.
func RenderHistoryTable(w io.Writer, entries []Entry) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "LABEL\tEXPRESSION\tRAN AT\tDURATION\tSTATUS\tMESSAGE")
	fmt.Fprintln(tw, "-----\t----------\t------\t--------\t------\t-------")
	for _, e := range entries {
		status := "OK"
		if !e.Success {
			status = "FAIL"
		}
		msg := e.Message
		if msg == "" {
			msg = "-"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\n",
			e.Label,
			e.Expression,
			e.RanAt.Format(time.RFC3339),
			e.Duration.Round(time.Millisecond).String(),
			status,
			msg,
		)
	}
	tw.Flush()
}

// RenderSuccessSummary writes a per-label success rate summary to w.
func RenderSuccessSummary(w io.Writer, store *Store) {
	seen := map[string]bool{}
	var labels []string
	for _, e := range store.All() {
		if !seen[e.Label] {
			seen[e.Label] = true
			labels = append(labels, e.Label)
		}
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "LABEL\tTOTAL\tSUCCESS RATE")
	fmt.Fprintln(tw, "-----\t-----\t------------")
	for _, label := range labels {
		entries := store.ByLabel(label)
		rate := store.SuccessRate(label)
		fmt.Fprintf(tw, "%s\t%d\t%.1f%%\n", label, len(entries), rate*100)
	}
	tw.Flush()
}
