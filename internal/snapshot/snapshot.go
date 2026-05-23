// Package snapshot provides functionality to save and load cron schedule
// snapshots to/from disk, enabling persistent state between crontrace runs.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Entry represents a single cron schedule entry stored in a snapshot.
type Entry struct {
	Label      string    `json:"label"`
	Expression string    `json:"expression"`
	NextRuns   []time.Time `json:"next_runs"`
	SavedAt    time.Time `json:"saved_at"`
}

// Snapshot holds a collection of cron entries captured at a point in time.
type Snapshot struct {
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	Entries   []Entry   `json:"entries"`
}

const currentVersion = 1

// New creates a new Snapshot with the provided entries and the current timestamp.
func New(entries []Entry) *Snapshot {
	return &Snapshot{
		Version:   currentVersion,
		CreatedAt: time.Now().UTC(),
		Entries:   entries,
	}
}

// Save serialises the snapshot to JSON and writes it to the given file path.
func Save(s *Snapshot, path string) error {
	if s == nil {
		return fmt.Errorf("snapshot: cannot save nil snapshot")
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal error: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("snapshot: write error: %w", err)
	}
	return nil
}

// Load reads a JSON snapshot from the given file path and deserialises it.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: read error: %w", err)
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal error: %w", err)
	}
	if s.Version != currentVersion {
		return nil, fmt.Errorf("snapshot: unsupported version %d (expected %d)", s.Version, currentVersion)
	}
	return &s, nil
}
