package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempIncidentStore(t *testing.T) (*IncidentStore, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "incidents.json")
	return NewIncidentStore(path), path
}

func TestIncidentStore_AppendAndLoad(t *testing.T) {
	store, _ := tempIncidentStore(t)
	inc := Incident{
		ID:        "inc-001",
		Host:      "192.168.1.1",
		Ports:     []int{22, 80},
		Kind:      "opened",
		Severity:  SeverityHigh,
		CreatedAt: time.Now().UTC(),
	}
	if err := store.Append(inc); err != nil {
		t.Fatalf("Append: %v", err)
	}
	results, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 incident, got %d", len(results))
	}
	if results[0].ID != "inc-001" {
		t.Errorf("expected ID inc-001, got %s", results[0].ID)
	}
}

func TestIncidentStore_Load_NoFile(t *testing.T) {
	store, _ := tempIncidentStore(t)
	results, err := store.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(results))
	}
}

func TestIncidentStore_Resolve(t *testing.T) {
	store, _ := tempIncidentStore(t)
	inc := Incident{
		ID:        "inc-002",
		Host:      "10.0.0.1",
		Ports:     []int{443},
		Kind:      "closed",
		Severity:  SeverityMedium,
		CreatedAt: time.Now().UTC(),
	}
	_ = store.Append(inc)
	if err := store.Resolve("inc-002"); err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	results, _ := store.Load()
	if results[0].ResolvedAt == nil {
		t.Error("expected ResolvedAt to be set")
	}
}

func TestIncidentStore_Append_InvalidDir(t *testing.T) {
	store := NewIncidentStore("/nonexistent/dir/incidents.json")
	inc := Incident{ID: "x", Host: "h", CreatedAt: time.Now().UTC()}
	if err := store.Append(inc); err == nil {
		t.Error("expected error for invalid path")
	}
}

func TestIncidentStore_MultipleEntries(t *testing.T) {
	store, _ := tempIncidentStore(t)
	for i, sev := range []IncidentSeverity{SeverityLow, SeverityMedium, SeverityHigh} {
		_ = store.Append(Incident{
			ID:        string(rune('a' + i)),
			Host:      "host",
			Severity:  sev,
			CreatedAt: time.Now().UTC(),
		})
	}
	results, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(results) != 3 {
		t.Errorf("expected 3 incidents, got %d", len(results))
	}
}

func TestIncidentStore_ResolveNoop(t *testing.T) {
	store, path := tempIncidentStore(t)
	// Resolve on empty store should not error and file should not be created
	if err := store.Resolve("missing-id"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(path); err == nil {
		// file was created; check it's an empty array
		results, _ := store.Load()
		if len(results) != 0 {
			t.Error("expected no incidents after noop resolve")
		}
	}
}
