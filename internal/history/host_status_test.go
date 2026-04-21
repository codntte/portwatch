package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempHostStatusStore(t *testing.T) (*HostStatusStore, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "host_status.json")
	return NewHostStatusStore(path), path
}

func TestHostStatusStore_UpsertAndGet(t *testing.T) {
	store, _ := tempHostStatusStore(t)
	now := time.Now().UTC().Truncate(time.Second)

	status := HostStatus{Host: "192.168.1.1", OpenPorts: []int{22, 80}, LastSeen: now, Up: true}
	if err := store.Upsert(status); err != nil {
		t.Fatalf("Upsert failed: %v", err)
	}

	got, err := store.Get("192.168.1.1")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got.Host != status.Host {
		t.Errorf("expected host %q, got %q", status.Host, got.Host)
	}
	if len(got.OpenPorts) != 2 {
		t.Errorf("expected 2 open ports, got %d", len(got.OpenPorts))
	}
	if !got.Up {
		t.Error("expected Up=true")
	}
}

func TestHostStatusStore_Load_NoFile(t *testing.T) {
	store, _ := tempHostStatusStore(t)
	records, err := store.LoadAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 0 {
		t.Errorf("expected empty slice, got %d records", len(records))
	}
}

func TestHostStatusStore_Get_NotFound(t *testing.T) {
	store, _ := tempHostStatusStore(t)
	_, err := store.Get("unknown-host")
	if err == nil {
		t.Error("expected error for missing host, got nil")
	}
}

func TestHostStatusStore_Upsert_Replaces(t *testing.T) {
	store, _ := tempHostStatusStore(t)
	now := time.Now().UTC().Truncate(time.Second)

	_ = store.Upsert(HostStatus{Host: "10.0.0.1", OpenPorts: []int{22}, LastSeen: now, Up: true})
	_ = store.Upsert(HostStatus{Host: "10.0.0.1", OpenPorts: []int{443, 8080}, LastSeen: now, Up: true})

	records, err := store.LoadAll()
	if err != nil {
		t.Fatalf("LoadAll failed: %v", err)
	}
	if len(records) != 1 {
		t.Errorf("expected 1 record after upsert, got %d", len(records))
	}
	if len(records[0].OpenPorts) != 2 {
		t.Errorf("expected updated ports, got %v", records[0].OpenPorts)
	}
}

func TestHostStatusStore_Upsert_InvalidDir(t *testing.T) {
	store := NewHostStatusStore("/nonexistent/dir/status.json")
	err := store.Upsert(HostStatus{Host: "x", Up: false})
	if err == nil {
		t.Error("expected error writing to invalid path")
	}
}

func TestHostStatusStore_LoadAll_Sorted(t *testing.T) {
	store, _ := tempHostStatusStore(t)
	now := time.Now().UTC()

	_ = store.Upsert(HostStatus{Host: "z-host", LastSeen: now})
	_ = store.Upsert(HostStatus{Host: "a-host", LastSeen: now})
	_ = store.Upsert(HostStatus{Host: "m-host", LastSeen: now})

	records, err := store.LoadAll()
	if err != nil {
		t.Fatalf("LoadAll failed: %v", err)
	}
	if records[0].Host != "a-host" || records[1].Host != "m-host" || records[2].Host != "z-host" {
		t.Errorf("records not sorted: %v", records)
	}
	os.Remove(store.path)
}
