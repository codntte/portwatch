package snapshot

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	orig := New("localhost", []int{22, 80, 443})
	if err := Save(path, orig); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if loaded.Host != orig.Host {
		t.Errorf("host: got %q, want %q", loaded.Host, orig.Host)
	}
	if len(loaded.OpenPorts) != len(orig.OpenPorts) {
		t.Errorf("ports length: got %d, want %d", len(loaded.OpenPorts), len(orig.OpenPorts))
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestCompare(t *testing.T) {
	prev := New("localhost", []int{22, 80, 443})
	curr := New("localhost", []int{22, 443, 8080})

	diff := Compare(prev, curr)

	sort.Ints(diff.Opened)
	sort.Ints(diff.Closed)

	if len(diff.Opened) != 1 || diff.Opened[0] != 8080 {
		t.Errorf("Opened: got %v, want [8080]", diff.Opened)
	}
	if len(diff.Closed) != 1 || diff.Closed[0] != 80 {
		t.Errorf("Closed: got %v, want [80]", diff.Closed)
	}
}

func TestCompare_NoDiff(t *testing.T) {
	prev := New("localhost", []int{22, 80})
	curr := New("localhost", []int{22, 80})

	diff := Compare(prev, curr)
	if len(diff.Opened) != 0 || len(diff.Closed) != 0 {
		t.Errorf("expected no diff, got opened=%v closed=%v", diff.Opened, diff.Closed)
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
