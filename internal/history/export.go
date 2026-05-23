package history

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// ExportEntry represents a single history record for export purposes.
type ExportEntry struct {
	Label      string    `json:"label"`
	Expression string    `json:"expression"`
	RanAt      time.Time `json:"ran_at"`
	Success    bool      `json:"success"`
	Message    string    `json:"message,omitempty"`
}

// ExportPayload wraps exported history records with metadata.
type ExportPayload struct {
	GeneratedAt time.Time     `json:"generated_at"`
	Count       int           `json:"count"`
	Entries     []ExportEntry `json:"entries"`
}

// ExportJSON writes all history records from the store as JSON to the given writer.
// Returns an error if marshalling fails or the store is nil.
func ExportJSON(s *Store, w io.Writer) error {
	if s == nil {
		return fmt.Errorf("history: store must not be nil")
	}

	records := s.All()
	entries := make([]ExportEntry, 0, len(records))
	for _, r := range records {
		entries = append(entries, ExportEntry{
			Label:      r.Label,
			Expression: r.Expression,
			RanAt:      r.RanAt,
			Success:    r.Success,
			Message:    r.Message,
		})
	}

	payload := ExportPayload{
		GeneratedAt: time.Now().UTC(),
		Count:       len(entries),
		Entries:     entries,
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(payload); err != nil {
		return fmt.Errorf("history: failed to encode JSON: %w", err)
	}
	return nil
}
