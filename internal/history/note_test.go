package history

import (
	"os"
	"path/filepath"
	"testing"
)

func tempNoteStore(t *testing.T) *NoteStore {
	t.Helper()
	dir := t.TempDir()
	return NewNoteStore(filepath.Join(dir, "notes.jsonl"))
}

func TestNoteStore_AddAndLoad(t *testing.T) {
	s := tempNoteStore(t)
	if err := s.Add("host-a", "first note"); err != nil {
		t.Fatalf("Add: %v", err)
	}
	if err := s.Add("host-b", "other note"); err != nil {
		t.Fatalf("Add: %v", err)
	}
	notes, err := s.Load("host-a")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(notes) != 1 {
		t.Fatalf("expected 1 note, got %d", len(notes))
	}
	if notes[0].Text != "first note" {
		t.Errorf("unexpected text: %q", notes[0].Text)
	}
}

func TestNoteStore_Load_AllHosts(t *testing.T) {
	s := tempNoteStore(t)
	_ = s.Add("host-a", "note 1")
	_ = s.Add("host-b", "note 2")
	notes, err := s.Load("")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(notes) != 2 {
		t.Errorf("expected 2 notes, got %d", len(notes))
	}
}

func TestNoteStore_Load_NoFile(t *testing.T) {
	s := tempNoteStore(t)
	notes, err := s.Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if notes != nil {
		t.Errorf("expected nil, got %v", notes)
	}
}

func TestNoteStore_Delete(t *testing.T) {
	s := tempNoteStore(t)
	_ = s.Add("host-a", "keep me not")
	_ = s.Add("host-b", "keep me")
	if err := s.Delete("host-a"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	notes, _ := s.Load("")
	if len(notes) != 1 || notes[0].Host != "host-b" {
		t.Errorf("unexpected notes after delete: %+v", notes)
	}
}

func TestNoteStore_Add_EmptyFields(t *testing.T) {
	s := tempNoteStore(t)
	if err := s.Add("", "text"); err == nil {
		t.Error("expected error for empty host")
	}
	if err := s.Add("host", ""); err == nil {
		t.Error("expected error for empty text")
	}
}

func TestNoteStore_Add_InvalidDir(t *testing.T) {
	s := NewNoteStore("/nonexistent/dir/notes.jsonl")
	if err := s.Add("host", "text"); err == nil {
		t.Error("expected error for invalid path")
	}
	_ = os.Remove("/nonexistent/dir/notes.jsonl")
}
