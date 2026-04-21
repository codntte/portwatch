package history

import (
	"os"
	"path/filepath"
	"testing"
)

func tempLabelStore(t *testing.T) *LabelStore {
	t.Helper()
	dir := t.TempDir()
	return NewLabelStore(filepath.Join(dir, "labels.json"))
}

func TestLabelStore_AddAndLoad(t *testing.T) {
	s := tempLabelStore(t)
	if err := s.Add("host1", "prod"); err != nil {
		t.Fatalf("Add: %v", err)
	}
	if err := s.Add("host2", "staging"); err != nil {
		t.Fatalf("Add: %v", err)
	}
	labels, err := s.Load("")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(labels) != 2 {
		t.Fatalf("expected 2 labels, got %d", len(labels))
	}
}

func TestLabelStore_Load_NoFile(t *testing.T) {
	s := tempLabelStore(t)
	labels, err := s.Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(labels) != 0 {
		t.Fatalf("expected empty, got %d", len(labels))
	}
}

func TestLabelStore_NoDuplicates(t *testing.T) {
	s := tempLabelStore(t)
	_ = s.Add("host1", "prod")
	_ = s.Add("host1", "prod")
	labels, _ := s.Load("")
	if len(labels) != 1 {
		t.Fatalf("expected 1 label (no duplicates), got %d", len(labels))
	}
}

func TestLabelStore_Delete(t *testing.T) {
	s := tempLabelStore(t)
	_ = s.Add("host1", "prod")
	_ = s.Add("host1", "critical")
	if err := s.Delete("host1", "prod"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	labels, _ := s.Load("host1")
	if len(labels) != 1 || labels[0].Name != "critical" {
		t.Fatalf("expected 1 remaining label 'critical', got %+v", labels)
	}
}

func TestLabelStore_FilterByHost(t *testing.T) {
	s := tempLabelStore(t)
	_ = s.Add("host1", "prod")
	_ = s.Add("host2", "staging")
	labels, _ := s.Load("host1")
	if len(labels) != 1 || labels[0].Host != "host1" {
		t.Fatalf("expected 1 label for host1, got %+v", labels)
	}
}

func TestLabelStore_Add_EmptyFields(t *testing.T) {
	s := tempLabelStore(t)
	if err := s.Add("", "prod"); err == nil {
		t.Fatal("expected error for empty host")
	}
	if err := s.Add("host1", ""); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestLabelStore_Add_InvalidDir(t *testing.T) {
	s := NewLabelStore("/nonexistent/dir/labels.json")
	if err := s.Add("host1", "prod"); err == nil {
		t.Fatal("expected error writing to invalid path")
	}
}

func TestLabelStore_Load_ByHost_NoMatch(t *testing.T) {
	s := tempLabelStore(t)
	_ = s.Add("host1", "prod")
	labels, err := s.Load("host99")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(labels) != 0 {
		t.Fatalf("expected 0 labels for unknown host, got %d", len(labels))
	}
}

func TestLabelStore_Persist(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "labels.json")
	s := NewLabelStore(path)
	_ = s.Add("host1", "prod")

	s2 := NewLabelStore(path)
	labels, err := s2.Load("")
	if err != nil {
		t.Fatalf("Load after reopen: %v", err)
	}
	if len(labels) != 1 {
		t.Fatalf("expected 1 persisted label, got %d", len(labels))
	}
	_ = os.Remove(path)
}
