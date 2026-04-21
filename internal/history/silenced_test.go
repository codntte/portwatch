package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempSilencedStore(t *testing.T) *SilencedStore {
	t.Helper()
	dir := t.TempDir()
	return NewSilencedStore(filepath.Join(dir, "silenced.json"))
}

func TestSilencedStore_AddAndLoad(t *testing.T) {
	s := tempSilencedStore(t)
	entry := SilencedEntry{
		Host:      "localhost",
		Port:      8080,
		Reason:    "maintenance",
		CreatedAt: time.Now(),
	}
	if err := s.Add(entry); err != nil {
		t.Fatalf("Add: %v", err)
	}
	entries, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Host != "localhost" || entries[0].Port != 8080 {
		t.Errorf("unexpected entry: %+v", entries[0])
	}
}

func TestSilencedStore_Load_NoFile(t *testing.T) {
	s := tempSilencedStore(t)
	entries, err := s.Load()
	if err != nil {
		t.Fatalf("expected nil error for missing file, got %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(entries))
	}
}

func TestSilencedStore_Add_Replaces(t *testing.T) {
	s := tempSilencedStore(t)
	s.Add(SilencedEntry{Host: "h1", Port: 22, Reason: "old", CreatedAt: time.Now()})
	s.Add(SilencedEntry{Host: "h1", Port: 22, Reason: "new", CreatedAt: time.Now()})
	entries, _ := s.Load()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry after replace, got %d", len(entries))
	}
	if entries[0].Reason != "new" {
		t.Errorf("expected reason 'new', got %q", entries[0].Reason)
	}
}

func TestSilencedStore_Delete(t *testing.T) {
	s := tempSilencedStore(t)
	s.Add(SilencedEntry{Host: "h1", Port: 443, Reason: "test", CreatedAt: time.Now()})
	if err := s.Delete("h1", 443); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	entries, _ := s.Load()
	if len(entries) != 0 {
		t.Errorf("expected 0 entries after delete, got %d", len(entries))
	}
}

func TestSilencedStore_IsSilenced(t *testing.T) {
	s := tempSilencedStore(t)
	s.Add(SilencedEntry{Host: "h1", Port: 80, Reason: "r", CreatedAt: time.Now()})
	if !s.IsSilenced("h1", 80) {
		t.Error("expected h1:80 to be silenced")
	}
	if s.IsSilenced("h1", 443) {
		t.Error("expected h1:443 to not be silenced")
	}
}

func TestSilencedStore_Expired(t *testing.T) {
	s := tempSilencedStore(t)
	past := time.Now().Add(-1 * time.Hour)
	s.Add(SilencedEntry{Host: "h1", Port: 9000, Reason: "expired", CreatedAt: past, ExpiresAt: past})
	entries, _ := s.Load()
	if len(entries) != 0 {
		t.Errorf("expected expired entry to be filtered, got %d entries", len(entries))
	}
}

func TestSilencedStore_Add_InvalidDir(t *testing.T) {
	s := NewSilencedStore(filepath.Join("/nonexistent", "silenced.json"))
	err := s.Add(SilencedEntry{Host: "h", Port: 1, CreatedAt: time.Now()})
	if err == nil {
		t.Error("expected error writing to invalid path")
	}
	os.Remove(filepath.Join("/nonexistent", "silenced.json"))
}
