package history

import (
	"os"
	"path/filepath"
	"testing"
)

func tempWatchlistStore(t *testing.T) (*WatchlistStore, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "watchlist.json")
	return NewWatchlistStore(path), path
}

func TestWatchlistStore_AddAndLoad(t *testing.T) {
	store, _ := tempWatchlistStore(t)

	if err := store.Add("localhost", 8080, "web"); err != nil {
		t.Fatalf("Add failed: %v", err)
	}
	if err := store.Add("192.168.1.1", 22, "ssh"); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	entries, err := store.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Host != "localhost" || entries[0].Port != 8080 {
		t.Errorf("unexpected first entry: %+v", entries[0])
	}
	if entries[1].Label != "ssh" {
		t.Errorf("expected label 'ssh', got %q", entries[1].Label)
	}
}

func TestWatchlistStore_NoDuplicates(t *testing.T) {
	store, _ := tempWatchlistStore(t)

	_ = store.Add("localhost", 9090, "")
	_ = store.Add("localhost", 9090, "dup")

	entries, _ := store.Load()
	if len(entries) != 1 {
		t.Errorf("expected 1 entry (no duplicates), got %d", len(entries))
	}
}

func TestWatchlistStore_Load_NoFile(t *testing.T) {
	store, _ := tempWatchlistStore(t)

	entries, err := store.Load()
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(entries))
	}
}

func TestWatchlistStore_Delete(t *testing.T) {
	store, _ := tempWatchlistStore(t)

	_ = store.Add("host-a", 80, "")
	_ = store.Add("host-b", 443, "")
	if err := store.Delete("host-a", 80); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	entries, _ := store.Load()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry after delete, got %d", len(entries))
	}
	if entries[0].Host != "host-b" {
		t.Errorf("expected host-b to remain, got %q", entries[0].Host)
	}
}

func TestWatchlistStore_Add_InvalidDir(t *testing.T) {
	store := NewWatchlistStore("/nonexistent/dir/watchlist.json")
	err := store.Add("host", 80, "")
	if err == nil {
		t.Error("expected error writing to invalid path")
	}
}

func TestWatchlistStore_AddedAt(t *testing.T) {
	store, _ := tempWatchlistStore(t)
	_ = store.Add("localhost", 3000, "dev")

	entries, _ := store.Load()
	if entries[0].AddedAt.IsZero() {
		t.Error("expected AddedAt to be set")
	}
}

func init() {
	// ensure os is imported for the invalid-dir test
	_ = os.DevNull
}
