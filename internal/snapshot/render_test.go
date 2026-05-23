package snapshot_test

import (
	"strings"
	"testing"
	"time"

	"github.com/user/crontrace/internal/snapshot"
)

func TestRenderSnapshotTable_NilSnapshot(t *testing.T) {
	var buf strings.Builder
	snapshot.RenderSnapshotTable(&buf, nil)
	if !strings.Contains(buf.String(), "no snapshot") {
		t.Errorf("expected 'no snapshot' message, got: %q", buf.String())
	}
}

func TestRenderSnapshotTable_HeaderPresent(t *testing.T) {
	s := snapshot.New([]snapshot.Entry{
		{Label: "daily", Expression: "0 0 * * *", NextRuns: []time.Time{time.Now().Add(time.Hour)}, SavedAt: time.Now()},
	})
	var buf strings.Builder
	snapshot.RenderSnapshotTable(&buf, s)
	out := buf.String()
	for _, col := range []string{"LABEL", "EXPRESSION", "NEXT RUN", "SAVED AT"} {
		if !strings.Contains(out, col) {
			t.Errorf("expected column %q in output", col)
		}
	}
}

func TestRenderSnapshotTable_EntryVisible(t *testing.T) {
	s := snapshot.New([]snapshot.Entry{
		{Label: "weekly", Expression: "0 0 * * 0", NextRuns: []time.Time{time.Now().Add(24 * time.Hour)}, SavedAt: time.Now()},
	})
	var buf strings.Builder
	snapshot.RenderSnapshotTable(&buf, s)
	out := buf.String()
	if !strings.Contains(out, "weekly") {
		t.Errorf("expected label 'weekly' in output, got: %q", out)
	}
	if !strings.Contains(out, "0 0 * * 0") {
		t.Errorf("expected expression in output, got: %q", out)
	}
}

func TestRenderSnapshotTable_DashWhenNoNextRuns(t *testing.T) {
	s := snapshot.New([]snapshot.Entry{
		{Label: "empty", Expression: "* * * * *", NextRuns: nil, SavedAt: time.Now()},
	})
	var buf strings.Builder
	snapshot.RenderSnapshotTable(&buf, s)
	if !strings.Contains(buf.String(), "-") {
		t.Error("expected dash for missing next run")
	}
}

func TestRenderSnapshotSummary_NilSnapshot(t *testing.T) {
	var buf strings.Builder
	snapshot.RenderSnapshotSummary(&buf, nil)
	if !strings.Contains(buf.String(), "no snapshot loaded") {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestRenderSnapshotSummary_CountVisible(t *testing.T) {
	s := snapshot.New([]snapshot.Entry{
		{Label: "a", Expression: "* * * * *"},
		{Label: "b", Expression: "0 * * * *"},
	})
	var buf strings.Builder
	snapshot.RenderSnapshotSummary(&buf, s)
	if !strings.Contains(buf.String(), "2 entries") {
		t.Errorf("expected '2 entries' in summary, got: %q", buf.String())
	}
}
