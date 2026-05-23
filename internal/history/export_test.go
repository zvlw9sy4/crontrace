package history

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestExportJSON_NilStore(t *testing.T) {
	var buf bytes.Buffer
	err := ExportJSON(nil, &buf)
	if err == nil {
		t.Fatal("expected error for nil store, got nil")
	}
	if !strings.Contains(err.Error(), "nil") {
		t.Errorf("expected error message to mention nil, got: %s", err.Error())
	}
}

func TestExportJSON_EmptyStore(t *testing.T) {
	s := NewStore()
	var buf bytes.Buffer
	if err := ExportJSON(s, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var payload ExportPayload
	if err := json.Unmarshal(buf.Bytes(), &payload); err != nil {
		t.Fatalf("failed to unmarshal output: %v", err)
	}
	if payload.Count != 0 {
		t.Errorf("expected count 0, got %d", payload.Count)
	}
	if len(payload.Entries) != 0 {
		t.Errorf("expected empty entries, got %d", len(payload.Entries))
	}
}

func TestExportJSON_WithEntries(t *testing.T) {
	s := NewStore()
	now := time.Now().UTC().Truncate(time.Second)

	_ = s.Record(Record{Label: "job-a", Expression: "0 * * * *", RanAt: now, Success: true, Message: ""})
	_ = s.Record(Record{Label: "job-b", Expression: "*/5 * * * *", RanAt: now.Add(time.Minute), Success: false, Message: "timeout"})

	var buf bytes.Buffer
	if err := ExportJSON(s, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var payload ExportPayload
	if err := json.Unmarshal(buf.Bytes(), &payload); err != nil {
		t.Fatalf("failed to unmarshal output: %v", err)
	}

	if payload.Count != 2 {
		t.Errorf("expected count 2, got %d", payload.Count)
	}
	if len(payload.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(payload.Entries))
	}

	found := false
	for _, e := range payload.Entries {
		if e.Label == "job-b" && e.Message == "timeout" && !e.Success {
			found = true
		}
	}
	if !found {
		t.Error("expected to find job-b entry with message 'timeout' and success=false")
	}
}

func TestExportJSON_GeneratedAtSet(t *testing.T) {
	s := NewStore()
	before := time.Now().UTC()
	var buf bytes.Buffer
	_ = ExportJSON(s, &buf)
	after := time.Now().UTC()

	var payload ExportPayload
	_ = json.Unmarshal(buf.Bytes(), &payload)

	if payload.GeneratedAt.Before(before) || payload.GeneratedAt.After(after) {
		t.Errorf("GeneratedAt %v is outside expected range [%v, %v]", payload.GeneratedAt, before, after)
	}
}
