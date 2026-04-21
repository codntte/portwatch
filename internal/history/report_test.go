package history

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeDiffEntry(t *testing.T, path, host string, ts time.Time) {
	t.Helper()
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	entry := DiffEntry{Host: host, Timestamp: ts, Opened: []int{80}, Closed: []int{}}
	if err := json.NewEncoder(f).Encode(entry); err != nil {
		t.Fatal(err)
	}
}

func writeStatEntry(t *testing.T, path, host string, ports []int) {
	t.Helper()
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	entry := StatsEntry{Host: host, OpenPorts: ports, ScannedAt: time.Now().UTC()}
	if err := json.NewEncoder(f).Encode(entry); err != nil {
		t.Fatal(err)
	}
}

func TestBuildReport_Basic(t *testing.T) {
	dir := t.TempDir()
	diffFile := filepath.Join(dir, "diffs.jsonl")
	statsFile := filepath.Join(dir, "stats.jsonl")

	now := time.Now().UTC()
	writeDiffEntry(t, diffFile, "host-a", now)
	writeDiffEntry(t, diffFile, "host-a", now.Add(time.Minute))
	writeStatEntry(t, statsFile, "host-a", []int{80, 443})

	report, err := BuildReport(diffFile, statsFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(report.Entries))
	}
	e := report.Entries[0]
	if e.Host != "host-a" {
		t.Errorf("expected host-a, got %s", e.Host)
	}
	if e.Changes != 2 {
		t.Errorf("expected 2 changes, got %d", e.Changes)
	}
	if len(e.OpenPorts) != 2 {
		t.Errorf("expected 2 open ports, got %d", len(e.OpenPorts))
	}
}

func TestBuildReport_NoFiles(t *testing.T) {
	dir := t.TempDir()
	report, err := BuildReport(
		filepath.Join(dir, "missing_diffs.jsonl"),
		filepath.Join(dir, "missing_stats.jsonl"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(report.Entries))
	}
}

func TestSaveReport(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "report.json")
	r := &Report{
		GeneratedAt: time.Now().UTC(),
		Entries: []ReportEntry{
			{Host: "host-b", OpenPorts: []int{22}, Changes: 1},
		},
	}
	if err := SaveReport(out, r); err != nil {
		t.Fatalf("save error: %v", err)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	var loaded Report
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(loaded.Entries) != 1 || loaded.Entries[0].Host != "host-b" {
		t.Errorf("unexpected loaded report: %+v", loaded)
	}
}
