package exporter_test

import (
	"strings"
	"testing"
	"time"

	"github.com/user/crontrace/internal/exporter"
)

func TestToCSV_HeaderPresent(t *testing.T) {
	p := exporter.Payload{
		Schedules: []exporter.ScheduleEntry{},
		Conflicts: []exporter.ConflictPair{},
	}
	out, err := exporter.ToCSV(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(out, "expression,description,next_run,conflict,conflict_with") {
		t.Errorf("expected CSV header as first line, got: %q", out)
	}
}

func TestToCSV_SingleEntry(t *testing.T) {
	now := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	p := exporter.Payload{
		Schedules: []exporter.ScheduleEntry{
			{
				Expression:  "0 9 * * 1",
				Description: "Every Monday at 09:00",
				NextRuns:    []time.Time{now},
			},
		},
		Conflicts: []exporter.ConflictPair{},
	}
	out, err := exporter.ToCSV(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "0 9 * * 1") {
		t.Errorf("expected expression in CSV output, got: %q", out)
	}
	if !strings.Contains(out, "false") {
		t.Errorf("expected conflict=false in output, got: %q", out)
	}
}

func TestToCSV_WithConflicts(t *testing.T) {
	now := time.Date(2024, 6, 1, 8, 0, 0, 0, time.UTC)
	p := exporter.Payload{
		Schedules: []exporter.ScheduleEntry{
			{Expression: "*/5 * * * *", Description: "Every 5 min", NextRuns: []time.Time{now}},
			{Expression: "*/5 * * * *", Description: "Also every 5 min", NextRuns: []time.Time{now}},
		},
		Conflicts: []exporter.ConflictPair{
			{A: "*/5 * * * *", B: "*/5 * * * *"},
		},
	}
	out, err := exporter.ToCSV(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "true") {
		t.Errorf("expected conflict=true in output, got: %q", out)
	}
}

func TestToCSV_EmptyNextRuns(t *testing.T) {
	p := exporter.Payload{
		Schedules: []exporter.ScheduleEntry{
			{Expression: "@yearly", Description: "Once a year", NextRuns: []time.Time{}},
		},
		Conflicts: []exporter.ConflictPair{},
	}
	out, err := exporter.ToCSV(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines (header + 1 row), got %d", len(lines))
	}
}
