package history_test

import (
	"os"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/history"
)

func TestAppendAndLoad(t *testing.T) {
	dir := t.TempDir()
	store := history.NewStore(dir)

	e1 := history.Entry{
		Timestamp: time.Now().UTC().Truncate(time.Second),
		Host:      "localhost",
		Opened:    []int{80, 443},
		Closed:    []int{},
	}
	e2 := history.Entry{
		Timestamp: time.Now().UTC().Truncate(time.Second),
		Host:      "localhost",
		Opened:    []int{},
		Closed:    []int{80},
	}

	if err := store.Append(e1); err != nil {
		t.Fatalf("Append e1: %v", err)
	}
	if err := store.Append(e2); err != nil {
		t.Fatalf("Append e2: %v", err)
	}

	entries, err := store.Load("localhost")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Host != "localhost" {
		t.Errorf("unexpected host: %s", entries[0].Host)
	}
	if len(entries[0].Opened) != 2 {
		t.Errorf("expected 2 opened ports, got %d", len(entries[0].Opened))
	}
}

func TestLoad_NoFile(t *testing.T) {
	store := history.NewStore(t.TempDir())
	entries, err := store.Load("ghost-host")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil entries for missing host")
	}
}

func TestAppend_InvalidDir(t *testing.T) {
	// Use a file as the directory to force an error.
	f, _ := os.CreateTemp("", "notadir")
	f.Close()
	defer os.Remove(f.Name())

	store := history.NewStore(f.Name())
	err := store.Append(history.Entry{Host: "h"})
	if err == nil {
		t.Fatal("expected error when dir is a file")
	}
}
