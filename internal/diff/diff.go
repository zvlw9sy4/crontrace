// Package diff provides utilities for comparing two sets of cron schedule
// entries and producing a human-readable change summary.
package diff

import (
	"fmt"
	"strings"

	"github.com/user/crontrace/internal/scheduler"
)

// ChangeKind describes the type of change between two schedule sets.
type ChangeKind string

const (
	Added    ChangeKind = "added"
	Removed  ChangeKind = "removed"
	Modified ChangeKind = "modified"
)

// Change represents a single difference between two schedule sets.
type Change struct {
	Kind     ChangeKind
	Label    string
	OldExpr  string
	NewExpr  string
}

// String returns a short human-readable description of the change.
func (c Change) String() string {
	switch c.Kind {
	case Added:
		return fmt.Sprintf("[+] %s  %s", c.Label, c.NewExpr)
	case Removed:
		return fmt.Sprintf("[-] %s  %s", c.Label, c.OldExpr)
	case Modified:
		return fmt.Sprintf("[~] %s  %s  →  %s", c.Label, c.OldExpr, c.NewExpr)
	}
	return ""
}

// Compare returns the list of changes between two slices of ScheduleEntry.
// Entries are matched by their Label field.
func Compare(before, after []scheduler.ScheduleEntry) []Change {
	oldMap := indexByLabel(before)
	newMap := indexByLabel(after)

	var changes []Change

	for label, oldEntry := range oldMap {
		if newEntry, ok := newMap[label]; ok {
			if strings.TrimSpace(oldEntry.Expression) != strings.TrimSpace(newEntry.Expression) {
				changes = append(changes, Change{
					Kind:    Modified,
					Label:   label,
					OldExpr: oldEntry.Expression,
					NewExpr: newEntry.Expression,
				})
			}
		} else {
			changes = append(changes, Change{
				Kind:    Removed,
				Label:   label,
				OldExpr: oldEntry.Expression,
			})
		}
	}

	for label, newEntry := range newMap {
		if _, ok := oldMap[label]; !ok {
			changes = append(changes, Change{
				Kind:    Added,
				Label:   label,
				NewExpr: newEntry.Expression,
			})
		}
	}

	return changes
}

func indexByLabel(entries []scheduler.ScheduleEntry) map[string]scheduler.ScheduleEntry {
	m := make(map[string]scheduler.ScheduleEntry, len(entries))
	for _, e := range entries {
		m[e.Label] = e
	}
	return m
}
