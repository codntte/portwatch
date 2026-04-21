package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempAlertLogStore(t *testing.T) (*AlertLogStore, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "alerts.jsonl")
	return NewAlertLogStore(path), path
}

func TestAlertLogStore_AppendAndLoad(t *testing.T) {
	store, _ := tempAlertLogStore(t)

	now := time.Now().UTC().Truncate(time.Second)
	entry := AlertEntry{
		Timestamp: now,
		Host:      "192.168.1.1",
		Opened:    []int{80, 443},
		Closed:    []int{22},
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
	if len(entries[0].Opened) != 2 {
		t.Errorf("expected 2 opened ports, got %d", len(entries[0].Opened))
	}
}

func TestAlertLogStore_Load_NoFile(t *testing.T) {
	store, _ := tempAlertLogStore(t)

	entries, err := store.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(entries))
	}
}

func TestAlertLogStore_MultipleEntries(t *testing.T) {
	store, _ := tempAlertLogStore(t)

	for i := 0; i < 3; i++ {
		err := store.Append(AlertEntry{
			Timestamp: time.Now().UTC(),
			Host:      "10.0.0.1",
			Opened:    []int{8080 + i},
		})
		if err != nil {
			t.Fatalf("Append[%d]: %v", i, err)
		}
	}

	entries, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(entries))
	}
}

func TestAlertLogStore_Append_InvalidDir(t *testing.T) {
	store := NewAlertLogStore("/nonexistent/dir/alerts.jsonl")
	err := store.Append(AlertEntry{
		Timestamp: time.Now().UTC(),
		Host:      "localhost",
	})
	if err == nil {
		t.Error("expected error for invalid directory, got nil")
	}
}

func TestAlertLogStore_Append_NoClosedOrOpened(t *testing.T) {
	store, path := tempAlertLogStore(t)

	err := store.Append(AlertEntry{
		Timestamp: time.Now().UTC(),
		Host:      "localhost",
	})
	if err != nil {
		t.Fatalf("Append: %v", err)
	}

	info, _ := os.Stat(path)
	if info.Size() == 0 {
		t.Error("expected non-empty file after append")
	}
}
