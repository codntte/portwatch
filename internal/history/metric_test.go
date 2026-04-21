package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempMetricStore(t *testing.T) (*MetricStore, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "metrics.jsonl")
	return NewMetricStore(path), path
}

func TestMetricStore_AppendAndLoad(t *testing.T) {
	store, _ := tempMetricStore(t)

	e := MetricEntry{Host: "host1", Name: "latency", Value: 12.5, Unit: "ms"}
	if err := store.Append(e); err != nil {
		t.Fatalf("Append: %v", err)
	}

	entries, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Value != 12.5 {
		t.Errorf("expected value 12.5, got %f", entries[0].Value)
	}
}

func TestMetricStore_Load_NoFile(t *testing.T) {
	store, _ := tempMetricStore(t)
	entries, err := store.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty, got %d", len(entries))
	}
}

func TestMetricStore_LoadByHost(t *testing.T) {
	store, _ := tempMetricStore(t)

	_ = store.Append(MetricEntry{Host: "a", Name: "cpu", Value: 50})
	_ = store.Append(MetricEntry{Host: "b", Name: "cpu", Value: 80})
	_ = store.Append(MetricEntry{Host: "a", Name: "mem", Value: 30})

	results, err := store.LoadByHost("a", "")
	if err != nil {
		t.Fatalf("LoadByHost: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 entries for host a, got %d", len(results))
	}

	results, err = store.LoadByHost("a", "cpu")
	if err != nil {
		t.Fatalf("LoadByHost with name: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 entry for host a / cpu, got %d", len(results))
	}
}

func TestMetricStore_Append_InvalidDir(t *testing.T) {
	store := NewMetricStore("/nonexistent/dir/metrics.jsonl")
	err := store.Append(MetricEntry{Host: "x", Name: "y", Value: 1})
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

func TestMetricStore_TimestampAutoSet(t *testing.T) {
	store, _ := tempMetricStore(t)
	before := time.Now().UTC()
	_ = store.Append(MetricEntry{Host: "h", Name: "n", Value: 1})
	after := time.Now().UTC()

	entries, _ := store.Load()
	if entries[0].Timestamp.Before(before) || entries[0].Timestamp.After(after) {
		t.Error("timestamp not set correctly")
	}
}

func TestMetricStore_MultipleEntries(t *testing.T) {
	store, _ := tempMetricStore(t)
	for i := 0; i < 5; i++ {
		_ = store.Append(MetricEntry{Host: "h", Name: "n", Value: float64(i)})
	}
	entries, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 5 {
		t.Errorf("expected 5 entries, got %d", len(entries))
	}
}

func TestMetricStore_PersistsAcrossInstances(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "metrics.jsonl")

	s1 := NewMetricStore(path)
	_ = s1.Append(MetricEntry{Host: "h", Name: "n", Value: 99})

	s2 := NewMetricStore(path)
	entries, err := s2.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 1 || entries[0].Value != 99 {
		t.Errorf("unexpected entries: %v", entries)
	}
	_ = os.Remove(path)
}
