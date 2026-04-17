package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestStats_Basic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")
	s := NewStore(path)

	now := time.Now()
	_ = s.Append(Entry{Host: "localhost", Timestamp: now, Opened: []int{80, 443}, Closed: []int{}})
	_ = s.Append(Entry{Host: "localhost", Timestamp: now.Add(time.Minute), Opened: []int{8080}, Closed: []int{80}})
	_ = s.Append(Entry{Host: "other", Timestamp: now, Opened: []int{22}, Closed: []int{}})

	stats, err := s.Stats("localhost")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stats) != 3 {
		t.Fatalf("expected 3 port stats, got %d", len(stats))
	}

	byPort := map[int]PortStat{}
	for _, st := range stats {
		byPort[st.Port] = st
	}

	if byPort[80].Opened != 1 || byPort[80].Closed != 1 {
		t.Errorf("port 80: expected opened=1 closed=1, got %+v", byPort[80])
	}
	if byPort[443].Opened != 1 || byPort[443].Closed != 0 {
		t.Errorf("port 443: expected opened=1 closed=0, got %+v", byPort[443])
	}
	if byPort[8080].Opened != 1 {
		t.Errorf("port 8080: expected opened=1, got %+v", byPort[8080])
	}
}

func TestStats_NoFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "missing.json")
	s := NewStore(path)

	stats, err := s.Stats("localhost")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stats) != 0 {
		t.Errorf("expected empty stats for missing file")
	}
}

func TestFormatStats(t *testing.T) {
	stats := []PortStat{
		{Port: 80, Opened: 3, Closed: 1},
		{Port: 443, Opened: 2, Closed: 0},
	}
	out := FormatStats(stats)
	if out == "" {
		t.Error("expected non-empty output")
	}
	_ = os.Stdout
}

func TestFormatStats_Empty(t *testing.T) {
	out := FormatStats(nil)
	if out != "no port activity recorded\n" {
		t.Errorf("unexpected output: %q", out)
	}
}
