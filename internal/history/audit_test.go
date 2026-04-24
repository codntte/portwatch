package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempAuditStore(t *testing.T) (*AuditStore, string) {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "audit.jsonl")
	return NewAuditStore(p), p
}

func TestAuditStore_AppendAndLoad(t *testing.T) {
	store, _ := tempAuditStore(t)

	entry := AuditEntry{
		Timestamp: time.Now().UTC(),
		Actor:     "admin",
		Action:    AuditActionAdd,
		Target:    "host:192.168.1.1",
		Detail:    "added to watchlist",
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
	if entries[0].Actor != "admin" {
		t.Errorf("expected actor 'admin', got %q", entries[0].Actor)
	}
	if entries[0].Action != AuditActionAdd {
		t.Errorf("expected action 'add', got %q", entries[0].Action)
	}
}

func TestAuditStore_Load_NoFile(t *testing.T) {
	store, _ := tempAuditStore(t)
	entries, err := store.Load()
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(entries))
	}
}

func TestAuditStore_MultipleEntries(t *testing.T) {
	store, _ := tempAuditStore(t)

	for i, action := range []AuditAction{AuditActionAdd, AuditActionUpdate, AuditActionDelete} {
		_ = i
		_ = store.Append(AuditEntry{Actor: "ci", Action: action, Target: "host:10.0.0.1"})
	}

	entries, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
}

func TestAuditStore_Append_InvalidDir(t *testing.T) {
	store := NewAuditStore("/nonexistent/dir/audit.jsonl")
	err := store.Append(AuditEntry{Actor: "x", Action: AuditActionAdd, Target: "y"})
	if err == nil {
		t.Fatal("expected error for invalid dir, got nil")
	}
}

func TestAuditStore_LoadByActor(t *testing.T) {
	store, _ := tempAuditStore(t)

	_ = store.Append(AuditEntry{Actor: "alice", Action: AuditActionAdd, Target: "host:1"})
	_ = store.Append(AuditEntry{Actor: "bob", Action: AuditActionDelete, Target: "host:2"})
	_ = store.Append(AuditEntry{Actor: "alice", Action: AuditActionUpdate, Target: "host:3"})

	entries, err := store.LoadByActor("alice")
	if err != nil {
		t.Fatalf("LoadByActor: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries for alice, got %d", len(entries))
	}
	for _, e := range entries {
		if e.Actor != "alice" {
			t.Errorf("unexpected actor %q", e.Actor)
		}
	}
}

func TestAuditStore_TimestampAutoSet(t *testing.T) {
	store, _ := tempAuditStore(t)
	before := time.Now().UTC()
	_ = store.Append(AuditEntry{Actor: "sys", Action: AuditActionAdd, Target: "host:5"})
	after := time.Now().UTC()

	entries, _ := store.Load()
	if len(entries) == 0 {
		t.Fatal("no entries loaded")
	}
	ts := entries[0].Timestamp
	if ts.Before(before) || ts.After(after) {
		t.Errorf("timestamp %v not in expected range [%v, %v]", ts, before, after)
	}
	_ = os.Remove("") // suppress unused import warning
}
