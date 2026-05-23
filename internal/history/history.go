// Package history tracks execution history entries for cron schedules,
// allowing users to record and query past run metadata.
package history

import (
	"errors"
	"time"
)

// Entry represents a single recorded execution of a cron schedule.
type Entry struct {
	Label      string
	Expression string
	RanAt      time.Time
	Duration   time.Duration
	Success    bool
	Message    string
}

// Store holds a collection of history entries.
type Store struct {
	entries []Entry
}

// NewStore returns an empty history Store.
func NewStore() *Store {
	return &Store{}
}

// Record appends a new Entry to the store.
func (s *Store) Record(e Entry) error {
	if e.Label == "" {
		return errors.New("history: label must not be empty")
	}
	if e.Expression == "" {
		return errors.New("history: expression must not be empty")
	}
	if e.RanAt.IsZero() {
		return errors.New("history: ran_at must not be zero")
	}
	s.entries = append(s.entries, e)
	return nil
}

// All returns a copy of all recorded entries.
func (s *Store) All() []Entry {
	out := make([]Entry, len(s.entries))
	copy(out, s.entries)
	return out
}

// ByLabel returns entries whose Label matches the given value.
func (s *Store) ByLabel(label string) []Entry {
	var out []Entry
	for _, e := range s.entries {
		if e.Label == label {
			out = append(out, e)
		}
	}
	return out
}

// SuccessRate returns the fraction of successful entries for a given label.
// Returns 0 if no entries exist for the label.
func (s *Store) SuccessRate(label string) float64 {
	entries := s.ByLabel(label)
	if len(entries) == 0 {
		return 0
	}
	var successes int
	for _, e := range entries {
		if e.Success {
			successes++
		}
	}
	return float64(successes) / float64(len(entries))
}
