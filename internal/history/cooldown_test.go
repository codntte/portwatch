package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempCooldownStore(t *testing.T) (*CooldownStore, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "cooldown.json")
	return NewCooldownStore(path), path
}

func TestCooldownStore_AddAndLoad(t *testing.T) {
	store, _ := tempCooldownStore(t)
	until := time.Now().UTC().Add(10 * time.Minute)

	if err := store.Add("host1", 8080, until); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	entries, err := store.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Host != "host1" || entries[0].Port != 8080 {
		t.Errorf("unexpected entry: %+v", entries[0])
	}
}

func TestCooldownStore_Load_NoFile(t *testing.T) {
	store, _ := tempCooldownStore(t)
	entries, err := store.Load()
	if err != nil {
		t.Fatalf("expected nil error for missing file, got %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(entries))
	}
}

func TestCooldownStore_IsActive(t *testing.T) {
	store, _ := tempCooldownStore(t)

	// Not active before adding.
	if store.IsActive("host1", 443) {
		t.Error("expected cooldown to be inactive before adding")
	}

	_ = store.Add("host1", 443, time.Now().UTC().Add(5*time.Minute))
	if !store.IsActive("host1", 443) {
		t.Error("expected cooldown to be active")
	}

	// Expired cooldown should not be active.
	_ = store.Add("host2", 80, time.Now().UTC().Add(-1*time.Minute))
	if store.IsActive("host2", 80) {
		t.Error("expected expired cooldown to be inactive")
	}
}

func TestCooldownStore_Add_Replaces(t *testing.T) {
	store, _ := tempCooldownStore(t)

	_ = store.Add("host1", 8080, time.Now().UTC().Add(1*time.Minute))
	_ = store.Add("host1", 8080, time.Now().UTC().Add(10*time.Minute))

	entries, _ := store.Load()
	if len(entries) != 1 {
		t.Errorf("expected 1 entry after replace, got %d", len(entries))
	}
}

func TestCooldownStore_Delete(t *testing.T) {
	store, _ := tempCooldownStore(t)

	_ = store.Add("host1", 22, time.Now().UTC().Add(5*time.Minute))
	_ = store.Add("host2", 22, time.Now().UTC().Add(5*time.Minute))
	_ = store.Delete("host1", 22)

	entries, _ := store.Load()
	if len(entries) != 1 {
		t.Errorf("expected 1 entry after delete, got %d", len(entries))
	}
	if entries[0].Host != "host2" {
		t.Errorf("unexpected remaining entry: %+v", entries[0])
	}
}

func TestCooldownStore_Add_InvalidDir(t *testing.T) {
	store := NewCooldownStore("/nonexistent/dir/cooldown.json")
	err := store.Add("host1", 80, time.Now().UTC().Add(time.Minute))
	if err == nil {
		t.Error("expected error writing to invalid path")
	}
	_ = os.Remove("/nonexistent/dir/cooldown.json")
}
