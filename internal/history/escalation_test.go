package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempEscalationStore(t *testing.T) (*EscalationStore, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "escalations.json")
	return NewEscalationStore(path), path
}

func TestEscalationStore_AppendAndLoad(t *testing.T) {
	store, _ := tempEscalationStore(t)
	now := time.Now().UTC().Truncate(time.Second)

	e := EscalationEntry{
		Host:      "192.168.1.1",
		Level:     EscalationCritical,
		Reason:    "port 443 closed unexpectedly",
		Triggered: now,
	}
	if err := store.Append(e); err != nil {
		t.Fatalf("Append: %v", err)
	}

	entries, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Host != "192.168.1.1" || entries[0].Level != EscalationCritical {
		t.Errorf("unexpected entry: %+v", entries[0])
	}
}

func TestEscalationStore_Load_NoFile(t *testing.T) {
	store, _ := tempEscalationStore(t)
	entries, err := store.Load()
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(entries))
	}
}

func TestEscalationStore_Resolve(t *testing.T) {
	store, _ := tempEscalationStore(t)
	now := time.Now().UTC().Truncate(time.Second)

	_ = store.Append(EscalationEntry{
		Host:      "10.0.0.1",
		Level:     EscalationWarning,
		Reason:    "high port churn",
		Triggered: now,
	})

	resolvedAt := now.Add(5 * time.Minute)
	if err := store.Resolve("10.0.0.1", resolvedAt); err != nil {
		t.Fatalf("Resolve: %v", err)
	}

	entries, _ := store.Load()
	if entries[0].Resolved == nil {
		t.Fatal("expected Resolved to be set")
	}
	if !entries[0].Resolved.Equal(resolvedAt) {
		t.Errorf("expected resolved at %v, got %v", resolvedAt, entries[0].Resolved)
	}
}

func TestEscalationStore_LoadByHost(t *testing.T) {
	store, _ := tempEscalationStore(t)
	now := time.Now().UTC()

	_ = store.Append(EscalationEntry{Host: "host-a", Level: EscalationInfo, Reason: "r1", Triggered: now})
	_ = store.Append(EscalationEntry{Host: "host-b", Level: EscalationCritical, Reason: "r2", Triggered: now})
	_ = store.Append(EscalationEntry{Host: "host-a", Level: EscalationWarning, Reason: "r3", Triggered: now})

	results, err := store.LoadByHost("host-a")
	if err != nil {
		t.Fatalf("LoadByHost: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 entries for host-a, got %d", len(results))
	}
}

func TestEscalationStore_Append_InvalidDir(t *testing.T) {
	store := NewEscalationStore("/nonexistent/dir/escalations.json")
	err := store.Append(EscalationEntry{
		Host:      "x",
		Level:     EscalationInfo,
		Reason:    "test",
		Triggered: time.Now(),
	})
	if err == nil {
		t.Fatal("expected error writing to invalid path")
	}
	_ = os.Remove("/nonexistent/dir/escalations.json")
}
