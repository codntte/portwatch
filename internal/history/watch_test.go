package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAppendEvent_AndLoad(t *testing.T) {
	dir := t.TempDir()
	s := &Store{path: filepath.Join(dir, "history.jsonl")}

	now := time.Now().UTC().Truncate(time.Second)
	events := []WatchEvent{
		{Timestamp: now, Host: "localhost", Opened: []int{80, 443}, Closed: nil},
		{Timestamp: now.Add(time.Minute), Host: "remote", Opened: nil, Closed: []int{22}},
	}

	for _, e := range events {
		if err := s.AppendEvent(e); err != nil {
			t.Fatalf("AppendEvent: %v", err)
		}
	}

	loaded, err := s.LoadEvents()
	if err != nil {
		t.Fatalf("LoadEvents: %v", err)
	}
	if len(loaded) != 2 {
		t.Fatalf("expected 2 events, got %d", len(loaded))
	}
	if loaded[0].Host != "localhost" {
		t.Errorf("expected host localhost, got %s", loaded[0].Host)
	}
	if len(loaded[0].Opened) != 2 {
		t.Errorf("expected 2 opened ports, got %d", len(loaded[0].Opened))
	}
	if loaded[1].Closed[0] != 22 {
		t.Errorf("expected closed port 22, got %d", loaded[1].Closed[0])
	}
}

func TestLoadEvents_NoFile(t *testing.T) {
	dir := t.TempDir()
	s := &Store{path: filepath.Join(dir, "missing.jsonl")}

	events, err := s.LoadEvents()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if events != nil {
		t.Errorf("expected nil, got %v", events)
	}
}

func TestAppendEvent_InvalidDir(t *testing.T) {
	s := &Store{path: "/nonexistent/dir/history.jsonl"}
	err := s.AppendEvent(WatchEvent{Timestamp: time.Now(), Host: "x"})
	if err == nil {
		t.Error("expected error for invalid dir")
	}
	_ = os.Remove("/nonexistent/dir/history.jsonl")
}
