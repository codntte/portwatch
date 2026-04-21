package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempMaintenanceStore(t *testing.T) (*MaintenanceStore, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "maintenance.jsonl")
	return NewMaintenanceStore(path), path
}

func TestMaintenanceStore_AddAndLoad(t *testing.T) {
	store, _ := tempMaintenanceStore(t)
	now := time.Now().UTC().Truncate(time.Second)
	w := MaintenanceWindow{
		ID:       "mw-1",
		Host:     "192.168.1.1",
		StartsAt: now.Add(-time.Minute),
		EndsAt:   now.Add(time.Hour),
		Reason:   "planned upgrade",
		CreatedAt: now,
	}
	if err := store.Add(w); err != nil {
		t.Fatalf("Add: %v", err)
	}
	windows, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(windows) != 1 {
		t.Fatalf("expected 1 window, got %d", len(windows))
	}
	if windows[0].ID != "mw-1" {
		t.Errorf("expected ID mw-1, got %s", windows[0].ID)
	}
}

func TestMaintenanceStore_Load_NoFile(t *testing.T) {
	store, _ := tempMaintenanceStore(t)
	windows, err := store.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(windows) != 0 {
		t.Errorf("expected empty slice, got %d", len(windows))
	}
}

func TestMaintenanceStore_Delete(t *testing.T) {
	store, _ := tempMaintenanceStore(t)
	now := time.Now().UTC()
	for _, id := range []string{"mw-1", "mw-2"} {
		_ = store.Add(MaintenanceWindow{ID: id, Host: "h1", StartsAt: now, EndsAt: now.Add(time.Hour)})
	}
	if err := store.Delete("mw-1"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	windows, _ := store.Load()
	if len(windows) != 1 || windows[0].ID != "mw-2" {
		t.Errorf("unexpected windows after delete: %+v", windows)
	}
}

func TestMaintenanceStore_ActiveFor(t *testing.T) {
	store, _ := tempMaintenanceStore(t)
	now := time.Now().UTC()
	_ = store.Add(MaintenanceWindow{ID: "active", Host: "h1", StartsAt: now.Add(-time.Minute), EndsAt: now.Add(time.Hour)})
	_ = store.Add(MaintenanceWindow{ID: "expired", Host: "h1", StartsAt: now.Add(-2 * time.Hour), EndsAt: now.Add(-time.Hour)})
	_ = store.Add(MaintenanceWindow{ID: "other", Host: "h2", StartsAt: now.Add(-time.Minute), EndsAt: now.Add(time.Hour)})

	active, err := store.ActiveFor("h1")
	if err != nil {
		t.Fatalf("ActiveFor: %v", err)
	}
	if len(active) != 1 || active[0].ID != "active" {
		t.Errorf("unexpected active windows: %+v", active)
	}
}

func TestMaintenanceStore_Add_InvalidDir(t *testing.T) {
	store := NewMaintenanceStore("/nonexistent/dir/maintenance.jsonl")
	err := store.Add(MaintenanceWindow{ID: "x"})
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

func TestMaintenanceWindow_IsActive(t *testing.T) {
	now := time.Now()
	w := MaintenanceWindow{StartsAt: now.Add(-time.Minute), EndsAt: now.Add(time.Hour)}
	if !w.IsActive() {
		t.Error("expected window to be active")
	}
	expired := MaintenanceWindow{StartsAt: now.Add(-2 * time.Hour), EndsAt: now.Add(-time.Hour)}
	if expired.IsActive() {
		t.Error("expected expired window to be inactive")
	}
}

func init() {
	// ensure file is cleaned up
	_ = os.Remove
	_ = filepath.Join
}
