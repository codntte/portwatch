package history

import (
	"os"
	"path/filepath"
	"testing"
)

func tempFingerprintStore(t *testing.T) (*FingerprintStore, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "fingerprints.jsonl")
	return NewFingerprintStore(path), path
}

func TestComputeFingerprint_Deterministic(t *testing.T) {
	a := ComputeFingerprint("localhost", []int{80, 443, 22})
	b := ComputeFingerprint("localhost", []int{443, 22, 80})
	if a != b {
		t.Errorf("expected same fingerprint for same ports in different order, got %s vs %s", a, b)
	}
}

func TestComputeFingerprint_DifferentHosts(t *testing.T) {
	a := ComputeFingerprint("host-a", []int{80})
	b := ComputeFingerprint("host-b", []int{80})
	if a == b {
		t.Error("expected different fingerprints for different hosts")
	}
}

func TestFingerprintStore_AppendAndLoad(t *testing.T) {
	store, _ := tempFingerprintStore(t)

	if err := store.Append("192.168.1.1", []int{22, 80}); err != nil {
		t.Fatalf("Append: %v", err)
	}
	if err := store.Append("192.168.1.2", []int{443}); err != nil {
		t.Fatalf("Append: %v", err)
	}

	entries, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Host != "192.168.1.1" {
		t.Errorf("unexpected host: %s", entries[0].Host)
	}
	if entries[0].Fingerprint == "" {
		t.Error("expected non-empty fingerprint")
	}
}

func TestFingerprintStore_Load_NoFile(t *testing.T) {
	store := NewFingerprintStore("/nonexistent/path/fp.jsonl")
	entries, err := store.Load()
	if err != nil {
		t.Fatalf("expected nil error for missing file, got %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil entries for missing file")
	}
}

func TestFingerprintStore_Append_InvalidDir(t *testing.T) {
	store := NewFingerprintStore("/nonexistent/dir/fp.jsonl")
	if err := store.Append("host", []int{80}); err == nil {
		t.Error("expected error for invalid directory")
	}
}

func TestFingerprintStore_LatestByHost(t *testing.T) {
	store, _ := tempFingerprintStore(t)

	_ = store.Append("host-a", []int{80})
	_ = store.Append("host-a", []int{80, 443})
	_ = store.Append("host-b", []int{22})

	latest, err := store.LatestByHost()
	if err != nil {
		t.Fatalf("LatestByHost: %v", err)
	}
	if len(latest) != 2 {
		t.Fatalf("expected 2 hosts, got %d", len(latest))
	}
	if len(latest["host-a"].Ports) != 2 {
		t.Errorf("expected latest entry for host-a to have 2 ports")
	}
}

func TestFingerprintStore_LatestByHost_NoFile(t *testing.T) {
	dir := t.TempDir()
	store := NewFingerprintStore(filepath.Join(dir, "missing.jsonl"))
	_ = os.Remove(filepath.Join(dir, "missing.jsonl"))

	latest, err := store.LatestByHost()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(latest) != 0 {
		t.Errorf("expected empty map for missing file")
	}
}
