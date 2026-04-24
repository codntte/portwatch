package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempHeartbeatStore(t *testing.T) *HeartbeatStore {
	t.Helper()
	dir := t.TempDir()
	return NewHeartbeatStore(filepath.Join(dir, "heartbeat.jsonl"))
}

func TestHeartbeatStore_AppendAndLoad(t *testing.T) {
	s := tempHeartbeatStore(t)
	now := time.Now().UTC().Truncate(time.Second)

	entry := HeartbeatEntry{Host: "host1", Timestamp: now, LatencyMs: 12, Alive: true}
	if err := s.Append(entry); err != nil {
		t.Fatalf("Append: %v", err)
	}

	entries, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Host != "host1" || entries[0].LatencyMs != 12 || !entries[0].Alive {
		t.Errorf("unexpected entry: %+v", entries[0])
	}
}

func TestHeartbeatStore_Load_NoFile(t *testing.T) {
	s := tempHeartbeatStore(t)
	entries, err := s.Load()
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(entries))
	}
}

func TestHeartbeatStore_LoadByHost(t *testing.T) {
	s := tempHeartbeatStore(t)
	now := time.Now().UTC()

	_ = s.Append(HeartbeatEntry{Host: "alpha", Timestamp: now, Alive: true})
	_ = s.Append(HeartbeatEntry{Host: "beta", Timestamp: now, Alive: false})
	_ = s.Append(HeartbeatEntry{Host: "alpha", Timestamp: now.Add(time.Second), Alive: true})

	entries, err := s.LoadByHost("alpha")
	if err != nil {
		t.Fatalf("LoadByHost: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries for alpha, got %d", len(entries))
	}
}

func TestHeartbeatStore_LastSeen(t *testing.T) {
	s := tempHeartbeatStore(t)
	now := time.Now().UTC().Truncate(time.Second)

	_ = s.Append(HeartbeatEntry{Host: "host1", Timestamp: now, LatencyMs: 5, Alive: true})
	_ = s.Append(HeartbeatEntry{Host: "host1", Timestamp: now.Add(time.Minute), LatencyMs: 8, Alive: true})

	last, err := s.LastSeen("host1")
	if err != nil {
		t.Fatalf("LastSeen: %v", err)
	}
	if last == nil {
		t.Fatal("expected a result, got nil")
	}
	if last.LatencyMs != 8 {
		t.Errorf("expected latency 8, got %d", last.LatencyMs)
	}
}

func TestHeartbeatStore_Append_InvalidDir(t *testing.T) {
	s := NewHeartbeatStore(filepath.Join(t.TempDir(), "no", "such", "dir", "hb.jsonl"))
	err := s.Append(HeartbeatEntry{Host: "x", Timestamp: time.Now(), Alive: true})
	if err == nil {
		t.Error("expected error for invalid path, got nil")
	}
}

func TestHeartbeatStore_LastSeen_NoFile(t *testing.T) {
	dir := t.TempDir()
	s := NewHeartbeatStore(filepath.Join(dir, "missing.jsonl"))
	_ = os.Remove(filepath.Join(dir, "missing.jsonl"))

	last, err := s.LastSeen("ghost")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if last != nil {
		t.Errorf("expected nil for unknown host, got %+v", last)
	}
}
