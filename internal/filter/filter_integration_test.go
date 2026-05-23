package filter_test

import (
	"testing"
	"time"

	"github.com/yourorg/crontrace/internal/filter"
)

// TestApply_CombinedCriteria verifies that multiple filter options are ANDed.
func TestApply_CombinedCriteria(t *testing.T) {
	now := baseTime()
	threshold := now.Add(30 * time.Minute)

	entries := []filter.Entry{
		// matches label but not time range
		{Expression: "0 * * * *", Label: "Health hourly", NextRuns: []time.Time{now.Add(1 * time.Hour)}},
		// matches time range but not label
		{Expression: "*/5 * * * *", Label: "DB ping", NextRuns: []time.Time{now.Add(5 * time.Minute)}},
		// matches both
		{Expression: "*/10 * * * *", Label: "Health check fast", NextRuns: []time.Time{now.Add(10 * time.Minute)}},
	}

	result := filter.Apply(entries, filter.Options{
		LabelContains: "health",
		ActiveBefore:  &threshold,
	})

	if len(result) != 1 {
		t.Fatalf("expected 1 combined match, got %d: %v", len(result), result)
	}
	if result[0].Label != "Health check fast" {
		t.Errorf("unexpected label: %s", result[0].Label)
	}
}

// TestApply_EmptyEntries ensures no panic on empty input.
func TestApply_EmptyEntries(t *testing.T) {
	result := filter.Apply(nil, filter.Options{LabelContains: "anything"})
	if result != nil && len(result) != 0 {
		t.Fatalf("expected empty result for nil input, got %v", result)
	}
}

// TestApply_EntryWithNoNextRuns is excluded by any time-range filter.
func TestApply_EntryWithNoNextRuns(t *testing.T) {
	now := baseTime()
	after := now.Add(-1 * time.Hour)

	entries := []filter.Entry{
		{Expression: "@yearly", Label: "Annual job", NextRuns: nil},
	}

	result := filter.Apply(entries, filter.Options{ActiveAfter: &after})
	if len(result) != 0 {
		t.Fatalf("entry with no next runs should be excluded by time filter, got %v", result)
	}
}
