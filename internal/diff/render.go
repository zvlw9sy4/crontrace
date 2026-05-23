package diff

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// RenderDiffTable writes a formatted table of changes to w.
// It returns the number of changes rendered and any write error.
func RenderDiffTable(w io.Writer, changes []Change) (int, error) {
	if len(changes) == 0 {
		_, err := fmt.Fprintln(w, "No changes detected.")
		return 0, err
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "KIND\tLABEL\tOLD EXPRESSION\tNEW EXPRESSION")
	fmt.Fprintln(tw, "----\t-----\t--------------\t--------------")

	for _, c := range changes {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n",
			kindLabel(c.Kind),
			c.Label,
			dashIfEmpty(c.OldExpr),
			dashIfEmpty(c.NewExpr),
		)
	}

	if err := tw.Flush(); err != nil {
		return 0, err
	}
	return len(changes), nil
}

func kindLabel(k ChangeKind) string {
	switch k {
	case Added:
		return "added"
	case Removed:
		return "removed"
	case Modified:
		return "modified"
	}
	return string(k)
}

func dashIfEmpty(s string) string {
	if s == "" {
		return "-"
	}
	return s
}
