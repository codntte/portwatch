package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempBookmarkStore(t *testing.T) *BookmarkStore {
	t.Helper()
	dir := t.TempDir()
	return NewBookmarkStore(filepath.Join(dir, "bookmarks.json"))
}

func TestBookmarkStore_AddAndLoad(t *testing.T) {
	s := tempBookmarkStore(t)
	b := Bookmark{
		Name:      "before-deploy",
		Host:      "192.168.1.1",
		CreatedAt: time.Now().UTC().Truncate(time.Second),
		Note:      "pre-release snapshot",
		Ports:     []int{22, 80, 443},
	}
	if err := s.Add(b); err != nil {
		t.Fatalf("Add: %v", err)
	}
	marks, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(marks) != 1 {
		t.Fatalf("expected 1 bookmark, got %d", len(marks))
	}
	if marks[0].Name != b.Name || marks[0].Host != b.Host {
		t.Errorf("unexpected bookmark: %+v", marks[0])
	}
}

func TestBookmarkStore_Load_NoFile(t *testing.T) {
	s := tempBookmarkStore(t)
	marks, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(marks) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(marks))
	}
}

func TestBookmarkStore_Add_Replaces(t *testing.T) {
	s := tempBookmarkStore(t)
	b := Bookmark{Name: "snap", Host: "10.0.0.1", Ports: []int{80}}
	_ = s.Add(b)
	b.Ports = []int{80, 443}
	_ = s.Add(b)
	marks, _ := s.Load()
	if len(marks) != 1 {
		t.Fatalf("expected 1 bookmark after replace, got %d", len(marks))
	}
	if len(marks[0].Ports) != 2 {
		t.Errorf("expected updated ports, got %v", marks[0].Ports)
	}
}

func TestBookmarkStore_Delete(t *testing.T) {
	s := tempBookmarkStore(t)
	_ = s.Add(Bookmark{Name: "a", Host: "h1", Ports: []int{22}})
	_ = s.Add(Bookmark{Name: "b", Host: "h1", Ports: []int{80}})
	if err := s.Delete("a", "h1"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	marks, _ := s.Load()
	if len(marks) != 1 || marks[0].Name != "b" {
		t.Errorf("expected only bookmark 'b', got %+v", marks)
	}
}

func TestBookmarkStore_Add_InvalidDir(t *testing.T) {
	s := NewBookmarkStore("/nonexistent/dir/bookmarks.json")
	err := s.Add(Bookmark{Name: "x", Host: "h", Ports: []int{}})
	if err == nil {
		t.Error("expected error writing to invalid path")
	}
}

func TestBookmarkStore_FilterByHost(t *testing.T) {
	s := tempBookmarkStore(t)
	_ = s.Add(Bookmark{Name: "a", Host: "h1", Ports: []int{22}})
	_ = s.Add(Bookmark{Name: "b", Host: "h2", Ports: []int{80}})
	_ = s.Add(Bookmark{Name: "c", Host: "h1", Ports: []int{443}})

	marks, _ := s.Load()
	var h1 []Bookmark
	for _, m := range marks {
		if m.Host == "h1" {
			h1 = append(h1, m)
		}
	}
	if len(h1) != 2 {
		t.Errorf("expected 2 bookmarks for h1, got %d", len(h1))
	}
	_ = os.Remove(s.path)
}
