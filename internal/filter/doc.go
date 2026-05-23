// Package filter provides filtering utilities for cron schedule entries.
//
// It allows callers to narrow down a list of [Entry] values by label substring,
// expression prefix, or time-range constraints on the computed next-run times.
//
// Basic usage:
//
//	import "github.com/yourorg/crontrace/internal/filter"
//
//	after := time.Now()
//	before := time.Now().Add(24 * time.Hour)
//	filtered := filter.Apply(entries, filter.Options{
//		LabelContains: "backup",
//		ActiveAfter:   &after,
//		ActiveBefore:  &before,
//	})
//
// All criteria are ANDed together; an empty Options struct passes every entry
// through unchanged.
package filter
