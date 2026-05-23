package exporter_test

import (
	"strings"
	"testing"
	"time"

	"github.com/user/crontrace/internal/exporter"
	"github.com/user/crontrace/internal/scheduler"
)

func makeRuns(n int) []scheduler.ScheduledRun {
	base := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	runs := make([]scheduler.ScheduledRun, n)
	for i := range runs {
		runs[i] = scheduler.ScheduledRun{Time: base.Add(time.Duration(i) * time.Hour)}
	}
	return runs
}

func TestToICAL_ContainsCalendarWrapper(t *testing.T) {
	runs := makeRuns(1)
	out, err := exporter.ToICAL("0 * * * *", runs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "BEGIN:VCALENDAR") {
		t.Error("expected BEGIN:VCALENDAR in output")
	}
	if !strings.Contains(out, "END:VCALENDAR") {
		t.Error("expected END:VCALENDAR in output")
	}
}

func TestToICAL_EventCountMatchesRuns(t *testing.T) {
	runs := makeRuns(3)
	out, err := exporter.ToICAL("*/5 * * * *", runs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	count := strings.Count(out, "BEGIN:VEVENT")
	if count != 3 {
		t.Errorf("expected 3 VEVENT blocks, got %d", count)
	}
}

func TestToICAL_ContainsExpression(t *testing.T) {
	expr := "30 8 * * 1"
	runs := makeRuns(1)
	out, err := exporter.ToICAL(expr, runs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, expr) {
		t.Errorf("expected expression %q in output", expr)
	}
}

func TestToICAL_EmptyExpressionReturnsError(t *testing.T) {
	_, err := exporter.ToICAL("", makeRuns(1))
	if err == nil {
		t.Error("expected error for empty expression, got nil")
	}
}

func TestToICAL_EmptyRunsProducesNoEvents(t *testing.T) {
	out, err := exporter.ToICAL("0 0 * * *", []scheduler.ScheduledRun{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "BEGIN:VEVENT") {
		t.Error("expected no VEVENT blocks for empty runs")
	}
}
