package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/crontrace/internal/snapshot"
)

func makeEntry(label, expr string) snapshot.Entry {
	return snapshot.Entry{
		Label:      label,
		Expression: expr,
		NextRuns:   []time.Time{time.Now().Add(time.Hour).UTC()},
		SavedAt:    time.Now().UTC(),
	}
}

func TestNew_SetsVersion(t *testing.T) {
	s := snapshot.New(nil)
	if s.Version != 1 {
		t.Errorf("expected version 1, got %d", s.Version)
	}
}

func TestNew_SetsCreatedAt(t *testing.T) {
	before := time.Now().UTC()
	s := snapshot.New(nil)
	after := time.Now().UTC()
	if s.CreatedAt.Before(before) || s.CreatedAt.After(after) {
		t.Error("CreatedAt is not within expected range")
	}
}

func TestSaveLoad_RoundTrip(t *testing.T) {
	entries := []snapshot.Entry{makeEntry("backup", "0 2 * * *")}
	s := snapshot.New(entries)

	tmp := filepath.Join(t.TempDir(), "snap.json")
	if err := snapshot.Save(s, tmp); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := snapshot.Load(tmp)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(loaded.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(loaded.Entries))
	}
	if loaded.Entries[0].Label != "backup" {
		t.Errorf("expected label 'backup', got %q", loaded.Entries[0].Label)
	}
}

func TestSave_NilSnapshot(t *testing.T) {
	err := snapshot.Save(nil, "/tmp/noop.json")
	if err == nil {
		t.Error("expected error for nil snapshot")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "bad.json")
	_ = os.WriteFile(tmp, []byte("not json"), 0o644)
	_, err := snapshot.Load(tmp)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
