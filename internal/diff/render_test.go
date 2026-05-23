package diff_test

import (
	"strings"
	"testing"

	"github.com/user/crontrace/internal/diff"
)

func TestRenderDiffTable_NoChanges(t *testing.T) {
	var sb strings.Builder
	n, err := diff.RenderDiffTable(&sb, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 changes rendered, got %d", n)
	}
	if !strings.Contains(sb.String(), "No changes") {
		t.Errorf("expected 'No changes' message, got: %s", sb.String())
	}
}

func TestRenderDiffTable_HeaderPresent(t *testing.T) {
	changes := []diff.Change{
		{Kind: diff.Added, Label: "job-x", NewExpr: "0 * * * *"},
	}
	var sb strings.Builder
	n, err := diff.RenderDiffTable(&sb, changes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1 change rendered, got %d", n)
	}
	out := sb.String()
	for _, col := range []string{"KIND", "LABEL", "OLD EXPRESSION", "NEW EXPRESSION"} {
		if !strings.Contains(out, col) {
			t.Errorf("output missing column header %q", col)
		}
	}
}

func TestRenderDiffTable_DashForEmptyExpr(t *testing.T) {
	changes := []diff.Change{
		{Kind: diff.Added, Label: "job-new", NewExpr: "*/10 * * * *"},
		{Kind: diff.Removed, Label: "job-old", OldExpr: "0 0 * * *"},
	}
	var sb strings.Builder
	_, err := diff.RenderDiffTable(&sb, changes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	// Added row should have '-' for old expression
	lines := strings.Split(out, "\n")
	addedLine := ""
	for _, l := range lines {
		if strings.Contains(l, "job-new") {
			addedLine = l
			break
		}
	}
	if addedLine == "" {
		t.Fatal("could not find added entry line")
	}
	if !strings.Contains(addedLine, "-") {
		t.Errorf("expected dash placeholder in added line: %s", addedLine)
	}
}

func TestRenderDiffTable_MultipleChanges(t *testing.T) {
	changes := []diff.Change{
		{Kind: diff.Added, Label: "a", NewExpr: "* * * * *"},
		{Kind: diff.Removed, Label: "b", OldExpr: "0 * * * *"},
		{Kind: diff.Modified, Label: "c", OldExpr: "0 0 * * *", NewExpr: "0 1 * * *"},
	}
	var sb strings.Builder
	n, err := diff.RenderDiffTable(&sb, changes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 3 {
		t.Errorf("expected 3 changes rendered, got %d", n)
	}
}
