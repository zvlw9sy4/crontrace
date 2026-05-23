package filter

import "time"

// Criteria holds the optional filtering parameters that can be applied
// to a collection of schedule entries.
type Criteria struct {
	// LabelContains filters entries whose label contains this substring
	// (case-insensitive). Empty string disables this filter.
	LabelContains string

	// ExpressionPrefix filters entries whose cron expression starts with
	// the given prefix. Empty string disables this filter.
	ExpressionPrefix string

	// After retains only next-run times strictly after this moment.
	// Zero value disables this filter.
	After time.Time

	// Before retains only next-run times strictly before this moment.
	// Zero value disables this filter.
	Before time.Time

	// OnlyConflicting, when true, retains only entries that have at least
	// one next-run time shared with another entry in the same set.
	OnlyConflicting bool

	// MaxNextRuns caps the number of next-run times kept per entry.
	// Values <= 0 impose no cap.
	MaxNextRuns int
}

// IsZero reports whether the Criteria is effectively empty (no filters set).
func (c Criteria) IsZero() bool {
	return c.LabelContains == "" &&
		c.ExpressionPrefix == "" &&
		c.After.IsZero() &&
		c.Before.IsZero() &&
		!c.OnlyConflicting &&
		c.MaxNextRuns <= 0
}
