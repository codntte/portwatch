package history

import (
	"bytes"
	"encoding/csv"
	"io"
	"strings"
	"testing"
	"time"
)

func TestExportCSV(t *testing.T) {
	dir := t.TempDir()
	s := NewStore(dir)
	host := "localhost"

	now := time.Now().UTC().Truncate(time.Second)
	entry := Entry{
		Timestamp: now,
		Opened:    []int{80, 443},
		Closed:    []int{8080},
	}
	if err := s.Append(host, entry); err != nil {
		t.Fatalf("append: %v", err)
	}

	var buf bytes.Buffer
	if err := s.ExportCSV(host, &buf); err != nil {
		t.Fatalf("export csv: %v", err)
	}

	r := csv.NewReader(&buf)
	records, err := r.ReadAll()
	if err != nil {
		t.Fatalf("read csv: %v", err)
	}

	// header + 2 opened + 1 closed
	if len(records) != 4 {
		t.Fatalf("expected 4 records, got %d", len(records))
	}
	if records[0][0] != "timestamp" || records[0][1] != "event" || records[0][2] != "port" {
		t.Errorf("unexpected header: %v", records[0])
	}
	events := map[string]bool{}
	for _, rec := range records[1:] {
		events[rec[1]+":"+rec[2]] = true
	}
	for _, want := range []string{"opened:80", "opened:443", "closed:8080"} {
		if !events[want] {
			t.Errorf("missing expected record %q", want)
		}
	}
}

func TestExportCSV_NoFile(t *testing.T) {
	dir := t.TempDir()
	s := NewStore(dir)

	var buf bytes.Buffer
	err := s.ExportCSV("ghost", &buf)
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}

	lines := strings.TrimSpace(buf.String())
	r := csv.NewReader(strings.NewReader(lines))
	records, _ := r.ReadAll()
	if len(records) != 1 {
		t.Errorf("expected only header row, got %d rows", len(records))
	}
	_ = io.Discard
}
