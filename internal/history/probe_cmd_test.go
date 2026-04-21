package history

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestPrintProbes_Basic(t *testing.T) {
	store, path := tempProbeStore(t)

	_ = store.Append(ProbeResult{
		Host:      "10.0.0.1",
		Timestamp: time.Now().UTC(),
		LatencyMs: 5,
		Success:   true,
	})
	_ = store.Append(ProbeResult{
		Host:      "10.0.0.2",
		Timestamp: time.Now().UTC(),
		LatencyMs: 0,
		Success:   false,
		Error:     "timeout",
	})

	var buf bytes.Buffer
	if err := printProbesTo(&buf, path, "", time.Time{}); err != nil {
		t.Fatalf("printProbesTo: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "10.0.0.1") {
		t.Errorf("expected host 10.0.0.1 in output")
	}
	if !strings.Contains(out, "timeout") {
		t.Errorf("expected error 'timeout' in output")
	}
	if !strings.Contains(out, "fail") {
		t.Errorf("expected status 'fail' in output")
	}
}

func TestPrintProbes_FilterHost(t *testing.T) {
	store, path := tempProbeStore(t)

	_ = store.Append(ProbeResult{Host: "alpha", Timestamp: time.Now().UTC(), Success: true})
	_ = store.Append(ProbeResult{Host: "beta", Timestamp: time.Now().UTC(), Success: true})

	var buf bytes.Buffer
	if err := printProbesTo(&buf, path, "alpha", time.Time{}); err != nil {
		t.Fatalf("printProbesTo: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "beta") {
		t.Errorf("expected beta to be filtered out")
	}
	if !strings.Contains(out, "alpha") {
		t.Errorf("expected alpha in output")
	}
}

func TestPrintProbes_NoResults(t *testing.T) {
	_, path := tempProbeStore(t)

	var buf bytes.Buffer
	if err := printProbesTo(&buf, path, "", time.Time{}); err != nil {
		t.Fatalf("printProbesTo: %v", err)
	}
	if !strings.Contains(buf.String(), "no probe results") {
		t.Errorf("expected 'no probe results' message")
	}
}

func TestPrintProbes_SinceFilter(t *testing.T) {
	store, path := tempProbeStore(t)

	old := time.Now().UTC().Add(-2 * time.Hour)
	recent := time.Now().UTC()

	_ = store.Append(ProbeResult{Host: "h1", Timestamp: old, Success: true})
	_ = store.Append(ProbeResult{Host: "h1", Timestamp: recent, Success: true})

	var buf bytes.Buffer
	cutoff := time.Now().UTC().Add(-30 * time.Minute)
	if err := printProbesTo(&buf, path, "", cutoff); err != nil {
		t.Fatalf("printProbesTo: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	// header + 1 data row
	if len(lines) != 2 {
		t.Errorf("expected 2 lines (header+1 result), got %d: %s", len(lines), buf.String())
	}
}
