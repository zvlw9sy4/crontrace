package history

import (
	"testing"
	"time"
)

var baseTime = time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

func makeEntry(label, expr string, success bool, offset time.Duration) Entry {
	return Entry{
		Label:      label,
		Expression: expr,
		RanAt:      baseTime.Add(offset),
		Duration:   time.Second,
		Success:    success,
	}
}

func TestRecord_Valid(t *testing.T) {
	s := NewStore()
	err := s.Record(makeEntry("job1", "*/5 * * * *", true, 0))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.All()) != 1 {
		t.Errorf("expected 1 entry, got %d", len(s.All()))
	}
}

func TestRecord_MissingLabel(t *testing.T) {
	s := NewStore()
	e := makeEntry("", "*/5 * * * *", true, 0)
	if err := s.Record(e); err == nil {
		t.Error("expected error for empty label")
	}
}

func TestRecord_MissingExpression(t *testing.T) {
	s := NewStore()
	e := makeEntry("job1", "", true, 0)
	if err := s.Record(e); err == nil {
		t.Error("expected error for empty expression")
	}
}

func TestRecord_ZeroRanAt(t *testing.T) {
	s := NewStore()
	e := Entry{Label: "job1", Expression: "*/5 * * * *", Success: true}
	if err := s.Record(e); err == nil {
		t.Error("expected error for zero ran_at")
	}
}

func TestByLabel(t *testing.T) {
	s := NewStore()
	_ = s.Record(makeEntry("alpha", "0 * * * *", true, 0))
	_ = s.Record(makeEntry("beta", "0 * * * *", false, time.Minute))
	_ = s.Record(makeEntry("alpha", "0 * * * *", false, 2*time.Minute))

	result := s.ByLabel("alpha")
	if len(result) != 2 {
		t.Errorf("expected 2 entries for alpha, got %d", len(result))
	}
}

func TestSuccessRate(t *testing.T) {
	s := NewStore()
	_ = s.Record(makeEntry("job", "*/1 * * * *", true, 0))
	_ = s.Record(makeEntry("job", "*/1 * * * *", true, time.Minute))
	_ = s.Record(makeEntry("job", "*/1 * * * *", false, 2*time.Minute))

	rate := s.SuccessRate("job")
	if rate < 0.66 || rate > 0.67 {
		t.Errorf("expected ~0.666, got %f", rate)
	}
}

func TestSuccessRate_NoEntries(t *testing.T) {
	s := NewStore()
	if s.SuccessRate("missing") != 0 {
		t.Error("expected 0 for unknown label")
	}
}
