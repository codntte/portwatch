package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempScheduleLogStore(t *testing.T) (*ScheduleLogStore, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "schedule.jsonl")
	return NewScheduleLogStore(path), path
}

func TestScheduleLogStore_AppendAndLoad(t *testing.T) {
	store, _ := tempScheduleLogStore(t)

	entry := ScheduleEntry{
		Timestamp: time.Now().UTC().Truncate(time.Second),
		Host:      "192.168.1.1",
		Duration:  "120ms",
		PortsOpen: 3,
	}
	if err := store.Append(entry); err != nil {
		t.Fatalf("Append: %v", err)
	}

	entries, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Host != "192.168.1.1" {
		t.Errorf("expected host 192.168.1.1, got %s", entries[0].Host)
	}
	if entries[0].PortsOpen != 3 {
		t.Errorf("expected 3 ports open, got %d", entries[0].PortsOpen)
	}
}

func TestScheduleLogStore_Load_NoFile(t *testing.T) {
	store, _ := tempScheduleLogStore(t)

	entries, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(entries))
	}
}

func TestScheduleLogStore_MultipleEntries(t *testing.T) {
	store, _ := tempScheduleLogStore(t)

	hosts := []string{"10.0.0.1", "10.0.0.2", "10.0.0.3"}
	for _, h := range hosts {
		if err := store.Append(ScheduleEntry{Host: h, Timestamp: time.Now().UTC()}); err != nil {
			t.Fatalf("Append(%s): %v", h, err)
		}
	}

	entries, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
}

func TestScheduleLogStore_Append_InvalidDir(t *testing.T) {
	store := NewScheduleLogStore("/nonexistent/dir/schedule.jsonl")
	err := store.Append(ScheduleEntry{Host: "localhost"})
	if err == nil {
		t.Error("expected error for invalid dir, got nil")
	}
}

func TestScheduleLogStore_WithError(t *testing.T) {
	store, _ := tempScheduleLogStore(t)

	entry := ScheduleEntry{
		Timestamp: time.Now().UTC(),
		Host:      "bad-host",
		Error:     "connection refused",
	}
	if err := store.Append(entry); err != nil {
		t.Fatalf("Append: %v", err)
	}

	entries, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if entries[0].Error != "connection refused" {
		t.Errorf("expected error field, got %q", entries[0].Error)
	}
}

func TestScheduleLogStore_Path(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sched.jsonl")
	store := NewScheduleLogStore(path)

	_ = store.Append(ScheduleEntry{Host: "h1", Timestamp: time.Now().UTC()})

	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file at %s, got error: %v", path, err)
	}
}
