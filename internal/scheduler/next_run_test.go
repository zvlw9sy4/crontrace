package scheduler_test

import (
	"testing"
	"time"

	"github.com/user/crontrace/internal/scheduler"
)

var baseTime = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func TestNextRuns_Count(t *testing.T) {
	runs, err := scheduler.NextRuns("* * * * *", baseTime, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(runs) != 5 {
		t.Errorf("expected 5 runs, got %d", len(runs))
	}
}

func TestNextRuns_Ordering(t *testing.T) {
	runs, err := scheduler.NextRuns("*/15 * * * *", baseTime, 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 1; i < len(runs); i++ {
		if !runs[i].After(runs[i-1]) {
			t.Errorf("runs not in ascending order at index %d", i)
		}
	}
}

func TestNextRuns_InvalidExpr(t *testing.T) {
	_, err := scheduler.NextRuns("invalid", baseTime, 3)
	if err == nil {
		t.Error("expected error for invalid expression, got nil")
	}
}

func TestDetectConflicts_Found(t *testing.T) {
	// Both run every hour — they share the same minute, so conflicts expected
	exprs := []string{"0 * * * *", "0 * * * *"}
	// Use distinct expressions that fire at the same time
	exprs = []string{"0 * * * *", "0 */1 * * *"}
	conflicts, err := scheduler.DetectConflicts(exprs, baseTime, 2*time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(conflicts) == 0 {
		t.Error("expected at least one conflict, got none")
	}
}

func TestDetectConflicts_None(t *testing.T) {
	// One runs at minute 0, the other at minute 30 — no overlap within 1 min window
	exprs := []string{"0 * * * *", "30 * * * *"}
	conflicts, err := scheduler.DetectConflicts(exprs, baseTime, 2*time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(conflicts) != 0 {
		t.Errorf("expected no conflicts, got %d", len(conflicts))
	}
}

func TestDetectConflicts_InvalidExpr(t *testing.T) {
	_, err := scheduler.DetectConflicts([]string{"bad expr"}, baseTime, time.Hour)
	if err == nil {
		t.Error("expected error for invalid expression")
	}
}
