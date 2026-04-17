package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestQuery_ByHost(t *testing.T) {
	dir := t.TempDir()
	s := NewStore(filepath.Join(dir, "history.json"))

	now := time.Now()
	_ = s.Append(Entry{Host: "host-a", Timestamp: now, Opened: []int{80}, Closed: []int{}})
	_ = s.Append(Entry{Host: "host-b", Timestamp: now, Opened: []int{443}, Closed: []int{}})

	results, err := s.Query(QueryOptions{Host: "host-a"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Host != "host-a" {
		t.Errorf("expected 1 host-a entry, got %+v", results)
	}
}

func TestQuery_SinceUntil(t *testing.T) {
	dir := t.TempDir()
	s := NewStore(filepath.Join(dir, "history.json"))

	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	_ = s.Append(Entry{Host: "h", Timestamp: base, Opened: []int{22}})
	_ = s.Append(Entry{Host: "h", Timestamp: base.Add(2 * time.Hour), Opened: []int{80}})
	_ = s.Append(Entry{Host: "h", Timestamp: base.Add(4 * time.Hour), Opened: []int{443}})

	results, err := s.Query(QueryOptions{
		Since: base.Add(1 * time.Hour),
		Until: base.Add(3 * time.Hour),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Opened[0] != 80 {
		t.Errorf("expected 1 result with port 80, got %+v", results)
	}
}

func TestQuery_Limit(t *testing.T) {
	dir := t.TempDir()
	s := NewStore(filepath.Join(dir, "history.json"))

	now := time.Now()
	for i := 0; i < 5; i++ {
		_ = s.Append(Entry{Host: "h", Timestamp: now, Opened: []int{i}})
	}

	results, err := s.Query(QueryOptions{Limit: 3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 3 {
		t.Errorf("expected 3 results, got %d", len(results))
	}
}

func TestQuery_NoFile(t *testing.T) {
	s := NewStore(filepath.Join(os.TempDir(), "nonexistent_query_history.json"))
	results, err := s.Query(QueryOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}
