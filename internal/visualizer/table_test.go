package visualizer_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/crontrace/internal/visualizer"
)

func TestRenderScheduleTable_Output(t *testing.T) {
	ref := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	rows := []visualizer.ScheduleRow{
		{
			Expression:  "0 * * * *",
			Description: "Every hour",
			NextRuns:    []time.Time{ref, ref.Add(time.Hour)},
		},
	}

	var buf bytes.Buffer
	visualizer.RenderScheduleTable(&buf, rows)
	out := buf.String()

	if !strings.Contains(out, "EXPRESSION") {
		t.Error("expected header EXPRESSION in output")
	}
	if !strings.Contains(out, "0 * * * *") {
		t.Error("expected expression in output")
	}
	if !strings.Contains(out, "Every hour") {
		t.Error("expected description in output")
	}
	if !strings.Contains(out, "2024-01-15 10:00:00") {
		t.Error("expected formatted time in output")
	}
}

func TestRenderConflictTable_NoConflicts(t *testing.T) {
	var buf bytes.Buffer
	visualizer.RenderConflictTable(&buf, []visualizer.ConflictRow{})
	out := buf.String()

	if !strings.Contains(out, "No conflicts") {
		t.Errorf("expected no-conflict message, got: %s", out)
	}
}

func TestRenderConflictTable_WithConflicts(t *testing.T) {
	rows := []visualizer.ConflictRow{
		{
			ExpressionA: "0 * * * *",
			ExpressionB: "0 */1 * * *",
			At:          time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC),
		},
	}

	var buf bytes.Buffer
	visualizer.RenderConflictTable(&buf, rows)
	out := buf.String()

	if !strings.Contains(out, "Conflicts Detected") {
		t.Error("expected conflict header")
	}
	if !strings.Contains(out, "0 * * * *") {
		t.Error("expected expression A in output")
	}
	if !strings.Contains(out, "2024-01-15 11:00:00") {
		t.Error("expected conflict time in output")
	}
}
