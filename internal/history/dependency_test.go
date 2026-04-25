package history

import (
	"os"
	"path/filepath"
	"testing"
)

func tempDependencyStore(t *testing.T) *DependencyStore {
	t.Helper()
	dir := t.TempDir()
	return NewDependencyStore(filepath.Join(dir, "dependencies.json"))
}

func TestDependencyStore_AddAndLoad(t *testing.T) {
	s := tempDependencyStore(t)
	dep := Dependency{
		ID:        "dep-1",
		Host:      "web-01",
		DependsOn: "db-01",
		Ports:     []int{5432},
		Note:      "primary db",
	}
	if err := s.Add(dep); err != nil {
		t.Fatalf("Add: %v", err)
	}
	deps, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(deps) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(deps))
	}
	if deps[0].ID != "dep-1" {
		t.Errorf("expected ID dep-1, got %s", deps[0].ID)
	}
	if deps[0].CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestDependencyStore_Load_NoFile(t *testing.T) {
	s := tempDependencyStore(t)
	deps, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(deps) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(deps))
	}
}

func TestDependencyStore_Delete(t *testing.T) {
	s := tempDependencyStore(t)
	_ = s.Add(Dependency{ID: "dep-1", Host: "a", DependsOn: "b"})
	_ = s.Add(Dependency{ID: "dep-2", Host: "c", DependsOn: "d"})
	if err := s.Delete("dep-1"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	deps, _ := s.Load()
	if len(deps) != 1 {
		t.Fatalf("expected 1 entry after delete, got %d", len(deps))
	}
	if deps[0].ID != "dep-2" {
		t.Errorf("expected dep-2, got %s", deps[0].ID)
	}
}

func TestDependencyStore_LoadByHost(t *testing.T) {
	s := tempDependencyStore(t)
	_ = s.Add(Dependency{ID: "1", Host: "web-01", DependsOn: "db-01"})
	_ = s.Add(Dependency{ID: "2", Host: "api-01", DependsOn: "web-01"})
	_ = s.Add(Dependency{ID: "3", Host: "api-01", DependsOn: "cache-01"})
	results, err := s.LoadByHost("web-01")
	if err != nil {
		t.Fatalf("LoadByHost: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results for web-01, got %d", len(results))
	}
}

func TestDependencyStore_Add_InvalidDir(t *testing.T) {
	s := NewDependencyStore("/nonexistent/dir/deps.json")
	err := s.Add(Dependency{ID: "x", Host: "a", DependsOn: "b"})
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

func TestDependencyStore_MultiplePorts(t *testing.T) {
	s := tempDependencyStore(t)
	dep := Dependency{
		ID:        "dep-ports",
		Host:      "svc-a",
		DependsOn: "svc-b",
		Ports:     []int{80, 443, 8080},
	}
	_ = s.Add(dep)
	deps, _ := s.Load()
	if len(deps[0].Ports) != 3 {
		t.Errorf("expected 3 ports, got %d", len(deps[0].Ports))
	}
	_ = os.Remove(s.path)
}
