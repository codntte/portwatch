package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempSnapshotDiffStore(t *testing.T) (*SnapshotDiffStore, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "snapshot_diff.jsonl")
	return NewSnapshotDiffStore(path), path
}

func TestSnapshotDiffStore_AppendAndLoad(t *testing.T) {
	store, _ := tempSnapshotDiffStore(t)
	entry := SnapshotDiffEntry{
		Timestamp: time.Now().UTC(),
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
		t.Errorf("host mismatch: got %s", entries[0].Host)
	}
	if len(entries[0].Opened) != 2 {
		t.Errorf("expected 2 opened ports, got %d", len(entries[0].Opened))
	}
}

func TestSnapshotDiffStore_Load_NoFile(t *testing.T) {
	store, _ := tempSnapshotDiffStore(t)
	entries, err := store.Load()
	if err != nil {
		t.Fatalf("expected nil error for missing file, got %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil entries, got %v", entries)
	}
}

func TestSnapshotDiffStore_MultipleEntries(t *testing.T) {
	store, _ := tempSnapshotDiffStore(t)
	for i := 0; i < 3; i++ {
		e := SnapshotDiffEntry{
			Host:   "10.0.0.1",
			Opened: []int{i + 1},
		}
		if err := store.Append(e); err != nil {
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

func TestSnapshotDiffStore_Append_InvalidDir(t *testing.T) {
	store := NewSnapshotDiffStore("/nonexistent/dir/snap.jsonl")
	err := store.Append(SnapshotDiffEntry{Host: "x"})
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

func TestSnapshotDiffStore_TimestampAutoSet(t *testing.T) {
	store, _ := tempSnapshotDiffStore(t)
	before := time.Now().UTC()
	if err := store.Append(SnapshotDiffEntry{Host: "auto"}); err != nil {
		t.Fatalf("Append: %v", err)
	}
	after := time.Now().UTC()
	entries, _ := store.Load()
	if entries[0].Timestamp.Before(before) || entries[0].Timestamp.After(after) {
		t.Errorf("auto timestamp out of range: %v", entries[0].Timestamp)
	}
}

func TestSnapshotDiffStore_EmptyPortSlices(t *testing.T) {
	store, _ := tempSnapshotDiffStore(t)
	_ = os.Remove(store.path) // ensure clean
	e := SnapshotDiffEntry{Host: "empty", Opened: []int{}, Closed: []int{}}
	if err := store.Append(e); err != nil {
		t.Fatalf("Append: %v", err)
	}
	entries, _ := store.Load()
	if len(entries[0].Opened) != 0 || len(entries[0].Closed) != 0 {
		t.Errorf("expected empty slices")
	}
}
