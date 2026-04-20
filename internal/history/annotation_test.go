package history

import (
	"os"
	"path/filepath"
	"testing"
)

func tempAnnotationStore(t *testing.T) *AnnotationStore {
	t.Helper()
	dir := t.TempDir()
	return NewAnnotationStore(filepath.Join(dir, "annotations.json"))
}

func TestAnnotationStore_AddAndLoad(t *testing.T) {
	s := tempAnnotationStore(t)

	if err := s.Add("id1", "localhost", "initial baseline set"); err != nil {
		t.Fatalf("Add: %v", err)
	}
	if err := s.Add("id2", "192.168.1.1", "port 22 opened intentionally"); err != nil {
		t.Fatalf("Add: %v", err)
	}

	entries, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Host != "localhost" {
		t.Errorf("expected host localhost, got %s", entries[0].Host)
	}
	if entries[1].Note != "port 22 opened intentionally" {
		t.Errorf("unexpected note: %s", entries[1].Note)
	}
}

func TestAnnotationStore_Load_NoFile(t *testing.T) {
	s := tempAnnotationStore(t)
	entries, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(entries))
	}
}

func TestAnnotationStore_Delete(t *testing.T) {
	s := tempAnnotationStore(t)

	_ = s.Add("aaa", "host1", "note one")
	_ = s.Add("bbb", "host2", "note two")

	if err := s.Delete("aaa"); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	entries, _ := s.Load()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry after delete, got %d", len(entries))
	}
	if entries[0].ID != "bbb" {
		t.Errorf("expected remaining entry id=bbb, got %s", entries[0].ID)
	}
}

func TestAnnotationStore_Add_EmptyFields(t *testing.T) {
	s := tempAnnotationStore(t)

	if err := s.Add("", "host", "note"); err == nil {
		t.Error("expected error for empty id")
	}
	if err := s.Add("id", "", "note"); err == nil {
		t.Error("expected error for empty host")
	}
	if err := s.Add("id", "host", ""); err == nil {
		t.Error("expected error for empty note")
	}
}

func TestAnnotationStore_Add_InvalidDir(t *testing.T) {
	s := NewAnnotationStore(filepath.Join("/nonexistent", "annotations.json"))
	if err := s.Add("x", "h", "n"); err == nil {
		t.Error("expected error writing to invalid path")
	}
	_ = os.Remove(filepath.Join("/nonexistent", "annotations.json"))
}
