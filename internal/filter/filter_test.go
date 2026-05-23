package filter_test

import (
	"testing"
	"time"

	"github.com/yourorg/crontrace/internal/filter"
)

func baseTime() time.Time {
	return time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
}

func sampleEntries() []filter.Entry {
	now := baseTime()
	return []filter.Entry{
		{Expression: "0 * * * *", Label: "Hourly backup", NextRuns: []time.Time{now.Add(1 * time.Hour), now.Add(2 * time.Hour)}},
		{Expression: "0 0 * * *", Label: "Daily report", NextRuns: []time.Time{now.Add(12 * time.Hour)}},
		{Expression: "*/5 * * * *", Label: "Health check", NextRuns: []time.Time{now.Add(5 * time.Minute), now.Add(10 * time.Minute)}},
	}
}

func TestApply_NoFilter(t *testing.T) {
	entries := sampleEntries()
	result := filter.Apply(entries, filter.Options{})
	if len(result) != len(entries) {
		t.Fatalf("expected %d entries, got %d", len(entries), len(result))
	}
}

func TestApply_LabelContains(t *testing.T) {
	result := filter.Apply(sampleEntries(), filter.Options{LabelContains: "backup"})
	if len(result) != 1 || result[0].Label != "Hourly backup" {
		t.Fatalf("expected 1 entry matching 'backup', got %v", result)
	}
}

func TestApply_LabelContains_CaseInsensitive(t *testing.T) {
	result := filter.Apply(sampleEntries(), filter.Options{LabelContains: "DAILY"})
	if len(result) != 1 || result[0].Label != "Daily report" {
		t.Fatalf("expected 1 entry matching 'DAILY', got %v", result)
	}
}

func TestApply_ExpressionPrefix(t *testing.T) {
	result := filter.Apply(sampleEntries(), filter.Options{ExpressionPrefix: "*/"})
	if len(result) != 1 || result[0].Expression != "*/5 * * * *" {
		t.Fatalf("expected 1 entry with prefix '*/', got %v", result)
	}
}

func TestApply_ActiveAfter(t *testing.T) {
	now := baseTime()
	threshold := now.Add(6 * time.Hour)
	result := filter.Apply(sampleEntries(), filter.Options{ActiveAfter: &threshold})
	// Only hourly backup (2h) and daily report (12h) have runs after 6h
	// hourly backup has run at 1h and 2h — neither after 6h, so excluded
	// daily report has run at 12h — included
	// health check has runs at 5m and 10m — excluded
	if len(result) != 1 || result[0].Label != "Daily report" {
		t.Fatalf("expected only 'Daily report', got %v", result)
	}
}

func TestApply_ActiveBefore(t *testing.T) {
	now := baseTime()
	threshold := now.Add(30 * time.Minute)
	result := filter.Apply(sampleEntries(), filter.Options{ActiveBefore: &threshold})
	if len(result) != 1 || result[0].Label != "Health check" {
		t.Fatalf("expected only 'Health check', got %v", result)
	}
}

func TestApply_NoMatch(t *testing.T) {
	result := filter.Apply(sampleEntries(), filter.Options{LabelContains: "nonexistent"})
	if len(result) != 0 {
		t.Fatalf("expected 0 results, got %d", len(result))
	}
}
