// Package exporter provides functionality to export cron schedule data
// in various formats for downstream consumption.
package exporter

import (
	"encoding/json"
	"fmt"
	"time"
)

// ScheduleEntry represents a single cron schedule with its next run times.
type ScheduleEntry struct {
	Expression string      `json:"expression"`
	Description string    `json:"description,omitempty"`
	NextRuns    []time.Time `json:"next_runs"`
}

// ConflictEntry represents a detected conflict between two cron schedules.
type ConflictEntry struct {
	ExpressionA string    `json:"expression_a"`
	ExpressionB string    `json:"expression_b"`
	ConflictAt  time.Time `json:"conflict_at"`
}

// ExportPayload is the top-level structure exported to JSON.
type ExportPayload struct {
	GeneratedAt time.Time       `json:"generated_at"`
	Schedules   []ScheduleEntry `json:"schedules"`
	Conflicts   []ConflictEntry `json:"conflicts"`
}

// ToJSON serializes an ExportPayload to a pretty-printed JSON byte slice.
func ToJSON(payload ExportPayload) ([]byte, error) {
	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("exporter: failed to marshal payload: %w", err)
	}
	return data, nil
}

// NewPayload constructs an ExportPayload with the current timestamp.
func NewPayload(schedules []ScheduleEntry, conflicts []ConflictEntry) ExportPayload {
	return ExportPayload{
		GeneratedAt: time.Now().UTC(),
		Schedules:   schedules,
		Conflicts:   conflicts,
	}
}
