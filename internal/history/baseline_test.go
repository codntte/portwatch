package history

import (
	"os"
	"path/filepath"
	"testing"
)

func tempBaselineStore(t *testing.T) *BaselineStore {
	t.Helper()
	dir := t.TempDir()
	return NewBaselineStore(filepath.Join(dir, "baselines.json"))
}

func TestBaselineStore_AddAndGet(t *testing.T) {
	s := tempBaselineStore(t)

	err := s.Add("initial", "localhost", []int{80, 443})
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	b, ok, err := s.Get("initial", "localhost")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if !ok {
		t.Fatal("expected baseline to exist")
	}
	if b.Name != "initial" || b.Host != "localhost" {
		t.Errorf("unexpected baseline: %+v", b)
	}
	if len(b.Ports) != 2 {
		t.Errorf("expected 2 ports, got %d", len(b.Ports))
	}
}

func TestBaselineStore_Get_NotFound(t *testing.T) {
	s := tempBaselineStore(t)

	_, ok, err := s.Get("missing", "localhost")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected baseline to be absent")
	}
}

func TestBaselineStore_Add_Replaces(t *testing.T) {
	s := tempBaselineStore(t)

	_ = s.Add("prod", "host1", []int{22, 80})
	_ = s.Add("prod", "host1", []int{443})

	b, ok, err := s.Get("prod", "host1")
	if err != nil || !ok {
		t.Fatalf("Get failed: err=%v ok=%v", err, ok)
	}
	if len(b.Ports) != 1 || b.Ports[0] != 443 {
		t.Errorf("expected updated ports, got %v", b.Ports)
	}
}

func TestBaselineStore_Delete(t *testing.T) {
	s := tempBaselineStore(t)

	_ = s.Add("base", "remotehost", []int{8080})
	if err := s.Delete("base", "remotehost"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, ok, _ := s.Get("base", "remotehost")
	if ok {
		t.Fatal("expected baseline to be deleted")
	}
}

func TestBaselineStore_List(t *testing.T) {
	s := tempBaselineStore(t)

	_ = s.Add("a", "h1", []int{80})
	_ = s.Add("b", "h2", []int{443})

	list, err := s.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 baselines, got %d", len(list))
	}
}

func TestBaselineStore_Load_NoFile(t *testing.T) {
	dir := t.TempDir()
	s := NewBaselineStore(filepath.Join(dir, "nonexistent.json"))

	list, err := s.List()
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("expected empty list, got %d entries", len(list))
	}
}

func TestBaselineStore_Add_InvalidDir(t *testing.T) {
	s := NewBaselineStore("/nonexistent/path/baselines.json")
	err := s.Add("x", "y", []int{1})
	if err == nil {
		t.Fatal("expected error writing to invalid path")
	}
	_ = os.Remove("/nonexistent/path/baselines.json")
}
