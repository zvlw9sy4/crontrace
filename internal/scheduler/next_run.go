// Package scheduler provides utilities for computing next-run times
// and detecting scheduling conflicts for cron expressions.
package scheduler

import (
	"time"

	"github.com/robfig/cron/v3"
)

// NextRuns returns the next n scheduled run times for the given cron expression,
// starting from the provided base time.
func NextRuns(expr string, base time.Time, n int) ([]time.Time, error) {
	schedule, err := cron.ParseStandard(expr)
	if err != nil {
		return nil, err
	}

	runs := make([]time.Time, 0, n)
	t := base
	for i := 0; i < n; i++ {
		t = schedule.Next(t)
		runs = append(runs, t)
	}
	return runs, nil
}

// ConflictWindow defines the minimum duration between two schedules
// for them to be considered non-conflicting.
const ConflictWindow = time.Minute

// Conflict represents a detected overlap between two cron schedules.
type Conflict struct {
	ExprA  string
	ExprB  string
	TimeA  time.Time
	TimeB  time.Time
	Delta  time.Duration
}

// DetectConflicts checks pairs of cron expressions for runs that fall within
// ConflictWindow of each other over the next horizon duration from base.
func DetectConflicts(exprs []string, base time.Time, horizon time.Duration) ([]Conflict, error) {
	type entry struct {
		expr string
		t    time.Time
	}

	var all []entry
	for _, expr := range exprs {
		schedule, err := cron.ParseStandard(expr)
		if err != nil {
			return nil, err
		}
		t := base
		for {
			t = schedule.Next(t)
			if t.After(base.Add(horizon)) {
				break
			}
			all = append(all, entry{expr: expr, t: t})
		}
	}

	var conflicts []Conflict
	for i := 0; i < len(all); i++ {
		for j := i + 1; j < len(all); j++ {
			if all[i].expr == all[j].expr {
				continue
			}
			delta := all[j].t.Sub(all[i].t)
			if delta < 0 {
				delta = -delta
			}
			if delta < ConflictWindow {
				conflicts = append(conflicts, Conflict{
					ExprA: all[i].expr,
					ExprB: all[j].expr,
					TimeA: all[i].t,
					TimeB: all[j].t,
					Delta: delta,
				})
			}
		}
	}
	return conflicts, nil
}
