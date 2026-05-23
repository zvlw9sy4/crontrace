// Package exporter provides functionality to export cron schedule data
// to various output formats including JSON and CSV.
package exporter

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"time"
)

// CSVRecord represents a single row in the CSV export.
type CSVRecord struct {
	Expression  string
	Description string
	NextRun     time.Time
	Conflict    bool
	ConflictWith string
}

// ToCSV serializes a Payload into CSV format and returns the result as a string.
// Columns: expression, description, next_run, conflict, conflict_with
func ToCSV(p Payload) (string, error) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	header := []string{"expression", "description", "next_run", "conflict", "conflict_with"}
	if err := w.Write(header); err != nil {
		return "", fmt.Errorf("csv: writing header: %w", err)
	}

	conflictMap := buildConflictMap(p.Conflicts)

	for _, entry := range p.Schedules {
		nextRun := ""
		if len(entry.NextRuns) > 0 {
			nextRun = entry.NextRuns[0].UTC().Format(time.RFC3339)
		}

		conflictWith := ""
		hasConflict := false
		if partner, ok := conflictMap[entry.Expression]; ok {
			hasConflict = true
			conflictWith = partner
		}

		row := []string{
			entry.Expression,
			entry.Description,
			nextRun,
			fmt.Sprintf("%t", hasConflict),
			conflictWith,
		}
		if err := w.Write(row); err != nil {
			return "", fmt.Errorf("csv: writing row for %q: %w", entry.Expression, err)
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return "", fmt.Errorf("csv: flushing writer: %w", err)
	}

	return buf.String(), nil
}

// buildConflictMap returns a map from expression -> conflicting expression
// using the conflict pairs stored in the payload.
func buildConflictMap(conflicts []ConflictPair) map[string]string {
	m := make(map[string]string, len(conflicts)*2)
	for _, c := range conflicts {
		m[c.A] = c.B
		m[c.B] = c.A
	}
	return m
}
