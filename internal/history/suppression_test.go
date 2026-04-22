package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempSuppressionStore(t *testing.T) (*SuppressionStore, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "suppression.json")
	return NewSuppressionStore(path), path
}

func TestSuppressionStore_AddAndLoad(t *testing.T) {
	store, _ := tempSuppressionStore(t)
	rule := SuppressionRule{
		Host:      "192.168.1.1",
		Port:      80,
		Reason:    "planned maintenance",
		CreatedAt: time.Now(),
	}
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

func TestSuppressionStore_Load_NoFile(t *testing.T) {
	store, _ := tempSuppressionStore(t)
	rules, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(rules) != 0 {
		t.Errorf("expected empty slice, got %d", len(rules))
	}
}

func TestSuppressionStore_Add_Replaces(t *testing.T) {
	store, _ := tempSuppressionStore(t)
	base := SuppressionRule{Host: "10.0.0.1", Port: 443, Reason: "old", CreatedAt: time.Now()}
	if err := store.Add(base); err != nil {
		t.Fatal(err)
	}
	updated := SuppressionRule{Host: "10.0.0.1", Port: 443, Reason: "new reason", CreatedAt: time.Now()}
	if err := store.Add(updated); err != nil {
		t.Fatal(err)
	}
	rules, _ := store.Load()
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule after replace, got %d", len(rules))
	}
	if rules[0].Reason != "new reason" {
		t.Errorf("expected updated reason, got %q", rules[0].Reason)
	}
}

func TestSuppressionStore_Delete(t *testing.T) {
	store, _ := tempSuppressionStore(t)
	rule := SuppressionRule{Host: "10.0.0.2", Port: 22, Reason: "temp", CreatedAt: time.Now()}
	_ = store.Add(rule)
	if err := store.Delete("10.0.0.2", 22); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	rules, _ := store.Load()
	if len(rules) != 0 {
		t.Errorf("expected 0 rules after delete, got %d", len(rules))
	}
}

func TestSuppressionStore_IsSuppressed(t *testing.T) {
	store, _ := tempSuppressionStore(t)
	rule := SuppressionRule{Host: "host1", Port: 8080, Reason: "test", CreatedAt: time.Now()}
	_ = store.Add(rule)

	ok, err := store.IsSuppressed("host1", 8080)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("expected host1:8080 to be suppressed")
	}

	ok, _ = store.IsSuppressed("host1", 9090)
	if ok {
		t.Error("expected host1:9090 to not be suppressed")
	}
}

func TestSuppressionStore_IsExpired(t *testing.T) {
	store, _ := tempSuppressionStore(t)
	expired := SuppressionRule{
		Host:      "host2",
		Port:      53,
		Reason:    "expired",
		CreatedAt: time.Now().Add(-2 * time.Hour),
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	_ = store.Add(expired)

	ok, err := store.IsSuppressed("host2", 53)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Error("expected expired rule to not suppress")
	}
}

func TestSuppressionStore_Add_InvalidDir(t *testing.T) {
	store := NewSuppressionStore("/nonexistent/dir/suppression.json")
	err := store.Add(SuppressionRule{Host: "h", Port: 1, CreatedAt: time.Now()})
	if err == nil {
		t.Error("expected error writing to invalid path")
	}
	_ = os.Remove("/nonexistent/dir/suppression.json")
}
