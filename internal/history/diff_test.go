package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAppendDiff_AndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "diffs.jsonl")

	now := time.Now().UTC().Truncate(time.Second)
	entry := DiffEntry{
		Timestamp: now,
		Host:      "192.168.1.1",
		Opened:    []int{80, 443},
		Closed:    []int{22},
	}

	if err := AppendDiff(path, entry); err != nil {
		t.Fatalf("AppendDiff: %v", err)
	}

	entries, err := LoadDiffs(path)
	if err != nil {
		t.Fatalf("LoadDiffs: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	got := entries[0]
	if got.Host != entry.Host {
		t.Errorf("host: got %q, want %q", got.Host, entry.Host)
	}
	if len(got.Opened) != 2 || got.Opened[0] != 80 {
		t.Errorf("opened ports mismatch: %v", got.Opened)
	}
	if len(got.Closed) != 1 || got.Closed[0] != 22 {
		t.Errorf("closed ports mismatch: %v", got.Closed)
	}
	// Verify timestamp round-trips correctly through JSON serialisation.
	if !got.Timestamp.Equal(now) {
		t.Errorf("timestamp: got %v, want %v", got.Timestamp, now)
	}
}

func TestLoadDiffs_NoFile(t *testing.T) {
	entries, err := LoadDiffs("/nonexistent/path/diffs.jsonl")
	if err != nil {
		t.Fatalf("expected nil error for missing file, got %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil entries, got %v", entries)
	}
}

func TestAppendDiff_MultipleEntries(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "diffs.jsonl")

	for i := 0; i < 3; i++ {
		e := DiffEntry{
			Timestamp: time.Now().UTC(),
			Host:      "host-1",
			Opened:    []int{i + 1},
			Closed:    []int{},
		}
		if err := AppendDiff(path, e); err != nil {
			t.Fatalf("AppendDiff[%d]: %v", i, err)
		}
	}

	entries, err := LoadDiffs(path)
	if err != nil {
		t.Fatalf("LoadDiffs: %v", err)
	}
	if len(entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(entries))
	}
	// Verify that opened port values are preserved in order.
	for i, e := range entries {
		if len(e.Opened) != 1 || e.Opened[0] != i+1 {
			t.Errorf("entry[%d] opened: got %v, want [%d]", i, e.Opened, i+1)
		}
	}
}

func TestAppendDiff_InvalidDir(t *testing.T) {
	err := AppendDiff("/nonexistent/dir/diffs.jsonl", DiffEntry{})
	if err == nil {
		t.Error("expected error for invalid directory, got nil")
	}
	_ = os.Remove("/nonexistent/dir/diffs.jsonl")
}
