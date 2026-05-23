package diff_test

import (
	"testing"
	"time"

	"github.com/user/crontrace/internal/diff"
	"github.com/user/crontrace/internal/scheduler"
)

func entry(label, expr string) scheduler.ScheduleEntry {
	return scheduler.ScheduleEntry{
		Label:      label,
		Expression: expr,
		NextRuns:   []time.Time{},
	}
}

func TestCompare_Added(t *testing.T) {
	before := []scheduler.ScheduleEntry{entry("job-a", "0 * * * *")}
	after := []scheduler.ScheduleEntry{
		entry("job-a", "0 * * * *"),
		entry("job-b", "*/5 * * * *"),
	}

	changes := diff.Compare(before, after)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Kind != diff.Added || changes[0].Label != "job-b" {
		t.Errorf("unexpected change: %+v", changes[0])
	}
}

func TestCompare_Removed(t *testing.T) {
	before := []scheduler.ScheduleEntry{
		entry("job-a", "0 * * * *"),
		entry("job-b", "*/5 * * * *"),
	}
	after := []scheduler.ScheduleEntry{entry("job-a", "0 * * * *")}

	changes := diff.Compare(before, after)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Kind != diff.Removed || changes[0].Label != "job-b" {
		t.Errorf("unexpected change: %+v", changes[0])
	}
}

func TestCompare_Modified(t *testing.T) {
	before := []scheduler.ScheduleEntry{entry("job-a", "0 * * * *")}
	after := []scheduler.ScheduleEntry{entry("job-a", "0 */2 * * *")}

	changes := diff.Compare(before, after)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	c := changes[0]
	if c.Kind != diff.Modified || c.OldExpr != "0 * * * *" || c.NewExpr != "0 */2 * * *" {
		t.Errorf("unexpected change: %+v", c)
	}
}

func TestCompare_NoChanges(t *testing.T) {
	entries := []scheduler.ScheduleEntry{entry("job-a", "0 * * * *")}
	changes := diff.Compare(entries, entries)
	if len(changes) != 0 {
		t.Errorf("expected no changes, got %d", len(changes))
	}
}

func TestChange_String(t *testing.T) {
	cases := []struct {
		c    diff.Change
		want string
	}{
		{diff.Change{Kind: diff.Added, Label: "j", NewExpr: "* * * * *"}, "[+] j  * * * * *"},
		{diff.Change{Kind: diff.Removed, Label: "j", OldExpr: "* * * * *"}, "[-] j  * * * * *"},
		{diff.Change{Kind: diff.Modified, Label: "j", OldExpr: "0 * * * *", NewExpr: "*/5 * * * *"}, "[~] j  0 * * * *  →  */5 * * * *"},
	}
	for _, tc := range cases {
		if got := tc.c.String(); got != tc.want {
			t.Errorf("String() = %q, want %q", got, tc.want)
		}
	}
}
