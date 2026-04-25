package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempRemediationStore(t *testing.T) (*RemediationStore, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "remediation.json")
	return NewRemediationStore(path), path
}

func makeRemediation(id, host string, port int, action string, status RemediationStatus) RemediationEntry {
	now := time.Now().UTC()
	return RemediationEntry{
		ID:        id,
		Host:      host,
		Port:      port,
		Action:    action,
		Status:    status,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func TestRemediationStore_AppendAndLoad(t *testing.T) {
	store, _ := tempRemediationStore(t)
	e := makeRemediation("r1", "host-a", 443, "restart-service", RemediationPending)
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
	if entries[0].ID != "r1" || entries[0].Host != "host-a" {
		t.Errorf("unexpected entry: %+v", entries[0])
	}
}

func TestRemediationStore_Load_NoFile(t *testing.T) {
	store, _ := tempRemediationStore(t)
	entries, err := store.Load()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty, got %d entries", len(entries))
	}
}

func TestRemediationStore_UpdateStatus(t *testing.T) {
	store, _ := tempRemediationStore(t)
	e := makeRemediation("r2", "host-b", 80, "block-port", RemediationPending)
	_ = store.Append(e)
	if err := store.UpdateStatus("r2", RemediationApplied, "done"); err != nil {
		t.Fatalf("UpdateStatus: %v", err)
	}
	entries, _ := store.Load()
	if entries[0].Status != RemediationApplied {
		t.Errorf("expected applied, got %s", entries[0].Status)
	}
	if entries[0].Note != "done" {
		t.Errorf("expected note 'done', got %s", entries[0].Note)
	}
}

func TestRemediationStore_Delete(t *testing.T) {
	store, _ := tempRemediationStore(t)
	_ = store.Append(makeRemediation("r3", "host-c", 22, "close-port", RemediationPending))
	_ = store.Append(makeRemediation("r4", "host-c", 8080, "alert", RemediationFailed))
	if err := store.Delete("r3"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	entries, _ := store.Load()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry after delete, got %d", len(entries))
	}
	if entries[0].ID != "r4" {
		t.Errorf("expected r4 to remain, got %s", entries[0].ID)
	}
}

func TestRemediationStore_Append_InvalidDir(t *testing.T) {
	store := NewRemediationStore("/nonexistent/dir/remediation.json")
	e := makeRemediation("r5", "host-d", 9090, "noop", RemediationSkipped)
	if err := store.Append(e); err == nil {
		t.Error("expected error for invalid dir, got nil")
	}
}

func TestRemediationStore_MultipleEntries(t *testing.T) {
	store, _ := tempRemediationStore(t)
	for i, action := range []string{"restart", "block", "notify"} {
		e := makeRemediation(
			string(rune('a'+i)),
			"host-e",
			1000+i,
			action,
			RemediationPending,
		)
		_ = store.Append(e)
	}
	entries, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(entries))
	}
}

func TestRemediationStore_InvalidPath_Load(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "remediation.json")
	if err := os.WriteFile(path, []byte("not-json"), 0644); err != nil {
		t.Fatal(err)
	}
	store := NewRemediationStore(path)
	_, err := store.Load()
	if err == nil {
		t.Error("expected parse error, got nil")
	}
}
