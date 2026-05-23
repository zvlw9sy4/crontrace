// Package filter provides utilities for filtering cron schedule entries
// based on various criteria such as time range, expression pattern, or tags.
package filter

import (
	"strings"
	"time"
)

// Entry represents a single cron schedule entry with metadata.
type Entry struct {
	Expression string
	Label      string
	NextRuns   []time.Time
}

// Options holds the filtering criteria applied to a list of entries.
type Options struct {
	// LabelContains filters entries whose label contains this substring (case-insensitive).
	LabelContains string
	// ActiveAfter keeps only entries with at least one next run after this time.
	ActiveAfter *time.Time
	// ActiveBefore keeps only entries with at least one next run before this time.
	ActiveBefore *time.Time
	// ExpressionPrefix keeps only entries whose expression starts with this prefix.
	ExpressionPrefix string
}

// Apply filters the provided entries according to the given Options and
// returns the subset that satisfies all non-zero criteria.
func Apply(entries []Entry, opts Options) []Entry {
	var result []Entry
	for _, e := range entries {
		if !matchesLabel(e, opts.LabelContains) {
			continue
		}
		if !matchesExpressionPrefix(e, opts.ExpressionPrefix) {
			continue
		}
		if !matchesTimeRange(e, opts.ActiveAfter, opts.ActiveBefore) {
			continue
		}
		result = append(result, e)
	}
	return result
}

func matchesLabel(e Entry, substr string) bool {
	if substr == "" {
		return true
	}
	return strings.Contains(strings.ToLower(e.Label), strings.ToLower(substr))
}

func matchesExpressionPrefix(e Entry, prefix string) bool {
	if prefix == "" {
		return true
	}
	return strings.HasPrefix(e.Expression, prefix)
}

func matchesTimeRange(e Entry, after, before *time.Time) bool {
	if after == nil && before == nil {
		return true
	}
	for _, t := range e.NextRuns {
		afterOK := after == nil || t.After(*after)
		beforeOK := before == nil || t.Before(*before)
		if afterOK && beforeOK {
			return true
		}
	}
	return false
}
