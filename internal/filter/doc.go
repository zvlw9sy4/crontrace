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
//
// Time-range filtering
//
// ActiveAfter and ActiveBefore are both optional. When only ActiveAfter is set,
// entries whose next-run time falls after that instant are included. When only
// ActiveBefore is set, entries whose next-run time falls before that instant are
// included. When both are set, only entries whose next-run time falls within the
// half-open interval [ActiveAfter, ActiveBefore) are included.
package filter
