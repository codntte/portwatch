package history

import (
	"path/filepath"
	"testing"
)

func TestTagStore_AddAndLoad(t *testing.T) {
	dir := t.TempDir()
	s := NewTagStore(filepath.Join(dir, "tags.json"))

	if err := s.Add("v1.0", "initial release"); err != nil {
		t.Fatal(err)
	}
	if err := s.Add("v1.1", ""); err != nil {
		t.Fatal(err)
	}

	tags, err := s.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(tags))
	}
	if tags[0].Name != "v1.0" || tags[0].Note != "initial release" {
		t.Errorf("unexpected tag: %+v", tags[0])
	}
}

func TestTagStore_Load_NoFile(t *testing.T) {
	dir := t.TempDir()
	s := NewTagStore(filepath.Join(dir, "tags.json"))
	tags, err := s.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) != 0 {
		t.Errorf("expected empty, got %d", len(tags))
	}
}

func TestTagStore_Delete(t *testing.T) {
	dir := t.TempDir()
	s := NewTagStore(filepath.Join(dir, "tags.json"))

	_ = s.Add("keep", "")
	_ = s.Add("remove", "")
	_ = s.Add("keep", "duplicate keep")

	if err := s.Delete("remove"); err != nil {
		t.Fatal(err)
	}

	tags, _ := s.Load()
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags after delete, got %d", len(tags))
	}
	for _, tg := range tags {
		if tg.Name == "remove" {
			t.Error("deleted tag still present")
		}
	}
}

func TestTagStore_Add_InvalidDir(t *testing.T) {
	s := NewTagStore("/nonexistent/dir/tags.json")
	if err := s.Add("x", ""); err == nil {
		t.Error("expected error for invalid path")
	}
}
