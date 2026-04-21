package history

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestPrintReport_Basic(t *testing.T) {
	dir := t.TempDir()
	diffFile := filepath.Join(dir, "diffs.jsonl")
	statsFile := filepath.Join(dir, "stats.jsonl")

	now := time.Now().UTC()
	writeDiffEntry(t, diffFile, "host-x", now)
	writeStatEntry(t, statsFile, "host-x", []int{22, 80})

	var buf bytes.Buffer
	if err := printReportTo(&buf, diffFile, statsFile); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "host-x") {
		t.Errorf("expected host-x in output, got:\n%s", out)
	}
	if !strings.Contains(out, "22,80") {
		t.Errorf("expected ports 22,80 in output, got:\n%s", out)
	}
	if !strings.Contains(out, "Report generated at:") {
		t.Errorf("expected header in output, got:\n%s", out)
	}
}

func TestPrintReport_NoData(t *testing.T) {
	dir := t.TempDir()
	var buf bytes.Buffer
	err := printReportTo(&buf,
		filepath.Join(dir, "no_diffs.jsonl"),
		filepath.Join(dir, "no_stats.jsonl"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No report data available.") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestFormatReportPorts_Empty(t *testing.T) {
	if got := formatReportPorts(nil); got != "-" {
		t.Errorf("expected '-', got %q", got)
	}
}

func TestFormatReportPorts_Values(t *testing.T) {
	got := formatReportPorts([]int{80, 443, 8080})
	if got != "80,443,8080" {
		t.Errorf("expected '80,443,8080', got %q", got)
	}
}
