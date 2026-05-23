package exporter_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/user/crontrace/internal/exporter"
)

func TestToJSON_ValidPayload(t *testing.T) {
	now := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	payload := exporter.ExportPayload{
		GeneratedAt: now,
		Schedules: []exporter.ScheduleEntry{
			{
				Expression:  "0 * * * *",
				Description: "Every hour",
				NextRuns:    []time.Time{now.Add(time.Hour)},
			},
		},
		Conflicts: []exporter.ConflictEntry{},
	}

	data, err := exporter.ToJSON(payload)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected non-empty JSON output")
	}

	var out map[string]interface{}
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if _, ok := out["generated_at"]; !ok {
		t.Error("expected 'generated_at' field in JSON output")
	}
	if _, ok := out["schedules"]; !ok {
		t.Error("expected 'schedules' field in JSON output")
	}
}

func TestToJSON_WithConflicts(t *testing.T) {
	now := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	payload := exporter.NewPayload(
		[]exporter.ScheduleEntry{
			{Expression: "*/5 * * * *", NextRuns: []time.Time{now}},
			{Expression: "*/15 * * * *", NextRuns: []time.Time{now}},
		},
		[]exporter.ConflictEntry{
			{ExpressionA: "*/5 * * * *", ExpressionB: "*/15 * * * *", ConflictAt: now},
		},
	)

	data, err := exporter.ToJSON(payload)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out map[string]interface{}
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	conflicts, ok := out["conflicts"].([]interface{})
	if !ok || len(conflicts) != 1 {
		t.Errorf("expected 1 conflict entry, got %v", out["conflicts"])
	}
}

func TestNewPayload_SetsGeneratedAt(t *testing.T) {
	before := time.Now().UTC()
	payload := exporter.NewPayload(nil, nil)
	after := time.Now().UTC()

	if payload.GeneratedAt.Before(before) || payload.GeneratedAt.After(after) {
		t.Errorf("GeneratedAt %v is outside expected range [%v, %v]", payload.GeneratedAt, before, after)
	}
}
