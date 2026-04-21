package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempPolicyStore(t *testing.T) *PolicyStore {
	t.Helper()
	dir := t.TempDir()
	return NewPolicyStore(filepath.Join(dir, "policies.json"))
}

func TestPolicyStore_AddAndLoad(t *testing.T) {
	s := tempPolicyStore(t)
	entry := PolicyEntry{
		Name:      "fast-scan",
		Host:      "192.168.1.1",
		PortRange: "1-1024",
		Interval:  5 * time.Minute,
		Enabled:   true,
	}
	if err := s.Add(entry); err != nil {
		t.Fatalf("Add: %v", err)
	}
	entries, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Name != "fast-scan" {
		t.Errorf("expected name fast-scan, got %s", entries[0].Name)
	}
}

func TestPolicyStore_Load_NoFile(t *testing.T) {
	s := tempPolicyStore(t)
	entries, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestPolicyStore_Add_Replaces(t *testing.T) {
	s := tempPolicyStore(t)
	for _, pr := range []string{"1-100", "1-9999"} {
		if err := s.Add(PolicyEntry{Name: "p", Host: "h", PortRange: pr, Enabled: true}); err != nil {
			t.Fatalf("Add: %v", err)
		}
	}
	entries, _ := s.Load()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry after replace, got %d", len(entries))
	}
	if entries[0].PortRange != "1-9999" {
		t.Errorf("expected updated port range, got %s", entries[0].PortRange)
	}
}

func TestPolicyStore_Delete(t *testing.T) {
	s := tempPolicyStore(t)
	_ = s.Add(PolicyEntry{Name: "a", Host: "h", PortRange: "1-80", Enabled: true})
	_ = s.Add(PolicyEntry{Name: "b", Host: "h", PortRange: "1-443", Enabled: false})
	if err := s.Delete("a"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	entries, _ := s.Load()
	if len(entries) != 1 || entries[0].Name != "b" {
		t.Errorf("expected only entry b after delete")
	}
}

func TestPolicyStore_Get(t *testing.T) {
	s := tempPolicyStore(t)
	_ = s.Add(PolicyEntry{Name: "scan-web", Host: "10.0.0.1", PortRange: "80-443", Enabled: true})
	e, found, err := s.Get("scan-web")
	if err != nil || !found {
		t.Fatalf("Get: err=%v found=%v", err, found)
	}
	if e.Host != "10.0.0.1" {
		t.Errorf("expected host 10.0.0.1, got %s", e.Host)
	}
}

func TestPolicyStore_Add_InvalidDir(t *testing.T) {
	s := NewPolicyStore("/nonexistent/dir/policies.json")
	err := s.Add(PolicyEntry{Name: "x", Host: "h", PortRange: "1-10", Enabled: true})
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

func TestPolicyStore_CreatedAt_Set(t *testing.T) {
	s := tempPolicyStore(t)
	before := time.Now()
	_ = s.Add(PolicyEntry{Name: "ts", Host: "h", PortRange: "22-22", Enabled: true})
	after := time.Now()
	entries, _ := s.Load()
	if entries[0].CreatedAt.Before(before) || entries[0].CreatedAt.After(after) {
		t.Errorf("CreatedAt not set correctly: %v", entries[0].CreatedAt)
	}
	_ = os.Remove(s.path)
}
