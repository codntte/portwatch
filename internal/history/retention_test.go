package history

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func writeEntries(t *testing.T, path string, entries []Entry) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	for _, e := range entries {
		if err := enc.Encode(e); err != nil {
			t.Fatal(err)
		}
	}
}

func TestRetentionPolicy_Apply_MaxAge(t *testing.T) {
	tmp := t.TempDir()
	path := tmp + "/history.jsonl"

	now := time.Now()
	entries := []Entry{
		{Host: "a", Timestamp: now.Add(-10 * 24 * time.Hour)},
		{Host: "b", Timestamp: now.Add(-1 * time.Hour)},
	}
	writeEntries(t, path, entries)

	p := RetentionPolicy{MaxAge: 7 * 24 * time.Hour, MaxEntries: 1000}
	if err := p.Apply(path); err != nil {
		t.Fatal(err)
	}

	store := NewStore(path)
	result, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 || result[0].Host != "b" {
		t.Errorf("expected 1 entry with host b, got %+v", result)
	}
}

func TestRetentionPolicy_Apply_MaxEntries(t *testing.T) {
	tmp := t.TempDir()
	path := tmp + "/history.jsonl"

	now := time.Now()
	var entries []Entry
	for i := 0; i < 5; i++ {
		entries = append(entries, Entry{Host: "h", Timestamp: now.Add(-time.Duration(i) * time.Minute)})
	}
	writeEntries(t, path, entries)

	p := RetentionPolicy{MaxAge: 24 * time.Hour, MaxEntries: 3}
	if err := p.Apply(path); err != nil {
		t.Fatal(err)
	}

	store := NewStore(path)
	result, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 3 {
		t.Errorf("expected 3 entries, got %d", len(result))
	}
}

func TestRetentionPolicy_Apply_NoFile(t *testing.T) {
	p := DefaultRetentionPolicy()
	if err := p.Apply("/tmp/portwatch_no_such_file_xyz.jsonl"); err != nil {
		t.Errorf("expected nil for missing file, got %v", err)
	}
}

func TestDefaultRetentionPolicy(t *testing.T) {
	p := DefaultRetentionPolicy()
	if p.MaxEntries != 1000 {
		t.Errorf("unexpected MaxEntries: %d", p.MaxEntries)
	}
	if p.MaxAge != 7*24*time.Hour {
		t.Errorf("unexpected MaxAge: %v", p.MaxAge)
	}
}
