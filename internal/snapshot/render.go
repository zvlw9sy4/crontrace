package snapshot

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

// RenderSnapshotTable writes a human-readable table of snapshot entries to w.
func RenderSnapshotTable(w io.Writer, s *Snapshot) {
	if s == nil {
		fmt.Fprintln(w, "(no snapshot)")
		return
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Snapshot v%d — created %s\n", s.Version, s.CreatedAt.Format(time.RFC3339))
	fmt.Fprintln(tw, "LABEL\tEXPRESSION\tNEXT RUN\tSAVED AT")
	fmt.Fprintln(tw, "-----\t----------\t--------\t--------")

	for _, e := range s.Entries {
		nextRun := "-"
		if len(e.NextRuns) > 0 {
			nextRun = e.NextRuns[0].Format(time.RFC3339)
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n",
			e.Label,
			e.Expression,
			nextRun,
			e.SavedAt.Format(time.RFC3339),
		)
	}
	tw.Flush()
}

// RenderSnapshotSummary writes a one-line summary of the snapshot to w.
func RenderSnapshotSummary(w io.Writer, s *Snapshot) {
	if s == nil {
		fmt.Fprintln(w, "no snapshot loaded")
		return
	}
	fmt.Fprintf(w, "snapshot: %d entries, created at %s\n",
		len(s.Entries),
		s.CreatedAt.Format(time.RFC3339),
	)
}
