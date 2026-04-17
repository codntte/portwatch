package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSummarize_Basic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")
	s := NewStore(path)

	now := time.Now()

	entries := []Entry{
		{Host: "host-a", Timestamp: now.Add(-2 * time.Hour), Opened: []int{80, 443}, Closed: []int{}},
		{Host: "host-a", Timestamp: now.Add(-1 * time.Hour), Opened: []int{}, Closed: []int{80}},
		{Host: "host-b", Timestamp: now, Opened: []int{22}, Closed: []int{8080}},
	}
	for _, e := range entries {
		if err := s.Append(e); err != nil {
			t.Fatalf("append: %v", err)
		}
	}

	summaries, err := s.Summarize(Query{})
	if err != nil {
		t.Fatalf("summarize: %v", err)
	}
	if len(summaries) != 2 {
		t.Fatalf("expected 2 summaries, got %d", len(summaries))
	}

	a := summaries[0]
	if a.Host != "host-a" {
		t.Errorf("expected host-a, got %s", a.Host)
	}
	if a.TotalEvents != 2 {
		t.Errorf("expected 2 events, got %d", a.TotalEvents)
	}
	if a.Opened != 2 {
		t.Errorf("expected 2 opened, got %d", a.Opened)
	}
	if a.Closed != 1 {
		t.Errorf("expected 1 closed, got %d", a.Closed)
	}

	b := summaries[1]
	if b.Host != "host-b" {
		t.Errorf("expected host-b, got %s", b.Host)
	}
	if b.TotalEvents != 1 {
		t.Errorf("expected 1 event, got %d", b.TotalEvents)
	}
}

func TestSummarize_NoFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.json")
	s := NewStore(path)
	summaries, err := s.Summarize(Query{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(summaries) != 0 {
		t.Errorf("expected empty summaries, got %d", len(summaries))
	}
	_ = os.Remove(path)
}
