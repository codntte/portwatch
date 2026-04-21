package history

import (
	"os"
	"path/filepath"
	"testing"
)

func tempThresholdStore(t *testing.T) (*ThresholdStore, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "thresholds.json")
	return NewThresholdStore(path), path
}

func TestThresholdStore_AddAndLoad(t *testing.T) {
	store, _ := tempThresholdStore(t)
	rule := ThresholdRule{Host: "192.168.1.1", Port: 80, MaxClosed: 3, Window: "1h"}
	if err := store.Add(rule); err != nil {
		t.Fatalf("Add: %v", err)
	}
	rules, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
	if rules[0].Host != "192.168.1.1" || rules[0].Port != 80 {
		t.Errorf("unexpected rule: %+v", rules[0])
	}
}

func TestThresholdStore_Load_NoFile(t *testing.T) {
	store, _ := tempThresholdStore(t)
	rules, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(rules) != 0 {
		t.Errorf("expected empty slice, got %d", len(rules))
	}
}

func TestThresholdStore_Add_Replaces(t *testing.T) {
	store, _ := tempThresholdStore(t)
	rule := ThresholdRule{Host: "10.0.0.1", Port: 443, MaxClosed: 2, Window: "30m"}
	_ = store.Add(rule)
	rule.MaxClosed = 5
	_ = store.Add(rule)
	rules, _ := store.Load()
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule after replace, got %d", len(rules))
	}
	if rules[0].MaxClosed != 5 {
		t.Errorf("expected MaxClosed=5, got %d", rules[0].MaxClosed)
	}
}

func TestThresholdStore_Delete(t *testing.T) {
	store, _ := tempThresholdStore(t)
	_ = store.Add(ThresholdRule{Host: "host1", Port: 22, MaxClosed: 1, Window: "1h"})
	_ = store.Add(ThresholdRule{Host: "host2", Port: 80, MaxClosed: 2, Window: "2h"})
	if err := store.Delete("host1", 22); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	rules, _ := store.Load()
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule after delete, got %d", len(rules))
	}
	if rules[0].Host != "host2" {
		t.Errorf("unexpected remaining rule: %+v", rules[0])
	}
}

func TestThresholdStore_Add_InvalidDir(t *testing.T) {
	store := NewThresholdStore("/nonexistent/dir/thresholds.json")
	err := store.Add(ThresholdRule{Host: "h", Port: 1, MaxClosed: 1, Window: "1h"})
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

func TestThresholdStore_CreatedAt(t *testing.T) {
	store, _ := tempThresholdStore(t)
	_ = store.Add(ThresholdRule{Host: "h", Port: 8080, MaxClosed: 3, Window: "6h"})
	rules, _ := store.Load()
	if rules[0].CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestThresholdStore_MultipleRules(t *testing.T) {
	store, _ := tempThresholdStore(t)
	for _, port := range []int{22, 80, 443} {
		_ = store.Add(ThresholdRule{Host: "myhost", Port: port, MaxClosed: 2, Window: "1h"})
	}
	rules, _ := store.Load()
	if len(rules) != 3 {
		t.Errorf("expected 3 rules, got %d", len(rules))
	}
	_ = os.Remove(filepath.Join(t.TempDir(), "thresholds.json"))
}
