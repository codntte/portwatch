package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempChangelogStore(t *testing.T) (*ChangelogStore, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "changelog.jsonl")
	store, err := NewChangelogStore(path)
	if err != nil {
		t.Fatalf("NewChangelogStore: %v", err)
	}
	return store, path
}

func TestChangelogStore_AppendAndLoad(t *testing.T) {
	store, _ := tempChangelogStore(t)

	entry := ChangelogEntry{
		Timestamp: time.Now().UTC().Truncate(time.Second),
		Host:      "192.168.1.1",
		Event:     "opened",
		Port:      443,
		Note:      "HTTPS port opened",
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

	got := entries[0]
	if got.Host != entry.Host {
		t.Errorf("Host: got %q, want %q", got.Host, entry.Host)
	}
	if got.Event != entry.Event {
		t.Errorf("Event: got %q, want %q", got.Event, entry.Event)
	}
	if got.Port != entry.Port {
		t.Errorf("Port: got %d, want %d", got.Port, entry.Port)
	}
	if got.Note != entry.Note {
		t.Errorf("Note: got %q, want %q", got.Note, entry.Note)
	}
}

func TestChangelogStore_Load_NoFile(t *testing.T) {
	dir := t.TempDir()
	store, err := NewChangelogStore(filepath.Join(dir, "missing.jsonl"))
	if err != nil {
		t.Fatalf("NewChangelogStore: %v", err)
	}

	entries, err := store.Load()
	if err != nil {
		t.Fatalf("Load on missing file: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestChangelogStore_MultipleEntries(t *testing.T) {
	store, _ := tempChangelogStore(t)

	now := time.Now().UTC().Truncate(time.Second)
	for i, tc := range []struct {
		host  string
		event string
		port  int
	}{
		{"10.0.0.1", "opened", 22},
		{"10.0.0.1", "closed", 80},
		{"10.0.0.2", "opened", 8080},
	} {
		if err := store.Append(ChangelogEntry{
			Timestamp: now.Add(time.Duration(i) * time.Second),
			Host:      tc.host,
			Event:     tc.event,
			Port:      tc.port,
		}); err != nil {
			t.Fatalf("Append[%d]: %v", i, err)
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

func TestChangelogStore_Append_InvalidDir(t *testing.T) {
	store, err := NewChangelogStore("/nonexistent/dir/changelog.jsonl")
	if err != nil {
		t.Fatalf("NewChangelogStore: %v", err)
	}

	err = store.Append(ChangelogEntry{
		Timestamp: time.Now().UTC(),
		Host:      "localhost",
		Event:     "opened",
		Port:      9090,
	})
	if err == nil {
		t.Error("expected error writing to invalid dir, got nil")
	}
}

func TestChangelogStore_LoadByHost(t *testing.T) {
	store, _ := tempChangelogStore(t)

	now := time.Now().UTC().Truncate(time.Second)
	entries := []ChangelogEntry{
		{Timestamp: now, Host: "host-a", Event: "opened", Port: 22},
		{Timestamp: now.Add(time.Second), Host: "host-b", Event: "opened", Port: 80},
		{Timestamp: now.Add(2 * time.Second), Host: "host-a", Event: "closed", Port: 443},
	}
	for _, e := range entries {
		if err := store.Append(e); err != nil {
			t.Fatalf("Append: %v", err)
		}
	}

	results, err := store.LoadByHost("host-a")
	if err != nil {
		t.Fatalf("LoadByHost: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 entries for host-a, got %d", len(results))
	}
	for _, r := range results {
		if r.Host != "host-a" {
			t.Errorf("unexpected host %q in results", r.Host)
		}
	}
}

func TestChangelogStore_Clear(t *testing.T) {
	store, path := tempChangelogStore(t)

	_ = store.Append(ChangelogEntry{
		Timestamp: time.Now().UTC(),
		Host:      "localhost",
		Event:     "opened",
		Port:      8443,
	})

	if err := store.Clear(); err != nil {
		t.Fatalf("Clear: %v", err)
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("expected file to be removed after Clear")
	}

	entries, err := store.Load()
	if err != nil {
		t.Fatalf("Load after Clear: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries after Clear, got %d", len(entries))
	}
}
