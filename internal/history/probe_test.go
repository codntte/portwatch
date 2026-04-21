package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempProbeStore(t *testing.T) (*ProbeStore, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "probes.jsonl")
	return NewProbeStore(path), path
}

func TestProbeStore_AppendAndLoad(t *testing.T) {
	store, _ := tempProbeStore(t)

	r := ProbeResult{
		Host:      "192.168.1.1",
		Timestamp: time.Now().UTC().Truncate(time.Second),
		LatencyMs: 12,
		Success:   true,
	}
	if err := store.Append(r); err != nil {
		t.Fatalf("Append: %v", err)
	}

	results, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Host != r.Host || results[0].LatencyMs != r.LatencyMs {
		t.Errorf("unexpected result: %+v", results[0])
	}
}

func TestProbeStore_Load_NoFile(t *testing.T) {
	store, _ := tempProbeStore(t)
	results, err := store.Load()
	if err != nil {
		t.Fatalf("Load on missing file: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(results))
	}
}

func TestProbeStore_LoadByHost(t *testing.T) {
	store, _ := tempProbeStore(t)

	hosts := []string{"host-a", "host-b", "host-a"}
	for _, h := range hosts {
		_ = store.Append(ProbeResult{Host: h, Timestamp: time.Now().UTC(), Success: true})
	}

	results, err := store.LoadByHost("host-a")
	if err != nil {
		t.Fatalf("LoadByHost: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results for host-a, got %d", len(results))
	}
}

func TestProbeStore_Append_InvalidDir(t *testing.T) {
	store := NewProbeStore("/nonexistent/dir/probes.jsonl")
	err := store.Append(ProbeResult{Host: "x", Timestamp: time.Now().UTC()})
	if err == nil {
		t.Fatal("expected error for invalid directory")
	}
}

func TestProbeStore_MultipleEntries(t *testing.T) {
	store, _ := tempProbeStore(t)

	for i := 0; i < 5; i++ {
		_ = store.Append(ProbeResult{
			Host:      "10.0.0.1",
			Timestamp: time.Now().UTC(),
			LatencyMs: int64(i * 10),
			Success:   i%2 == 0,
		})
	}

	results, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(results) != 5 {
		t.Errorf("expected 5 entries, got %d", len(results))
	}
}

func TestProbeStore_ErrorEntry(t *testing.T) {
	store, _ := tempProbeStore(t)

	r := ProbeResult{
		Host:      "unreachable.host",
		Timestamp: time.Now().UTC(),
		Success:   false,
		Error:     "connection refused",
	}
	if err := store.Append(r); err != nil {
		t.Fatalf("Append: %v", err)
	}

	results, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if results[0].Error != "connection refused" {
		t.Errorf("expected error field, got %q", results[0].Error)
	}
	_ = os.Remove("") // suppress unused import warning
}
