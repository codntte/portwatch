package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempAcknowledgmentStore(t *testing.T) *AcknowledgmentStore {
	t.Helper()
	dir := t.TempDir()
	return NewAcknowledgmentStore(filepath.Join(dir, "acks.json"))
}

func TestAcknowledgmentStore_AppendAndLoad(t *testing.T) {
	s := tempAcknowledgmentStore(t)
	a := Acknowledgment{
		ID:      "ack-1",
		Host:    "192.168.1.1",
		Port:    443,
		AckedBy: "alice",
		Comment: "expected change",
		AckedAt: time.Now().UTC().Truncate(time.Second),
	}
	if err := s.Append(a); err != nil {
		t.Fatalf("Append: %v", err)
	}
	entries, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Host != a.Host || entries[0].Port != a.Port {
		t.Errorf("entry mismatch: got %+v", entries[0])
	}
}

func TestAcknowledgmentStore_Load_NoFile(t *testing.T) {
	s := tempAcknowledgmentStore(t)
	entries, err := s.Load()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(entries))
	}
}

func TestAcknowledgmentStore_LoadByHost(t *testing.T) {
	s := tempAcknowledgmentStore(t)
	_ = s.Append(Acknowledgment{ID: "a1", Host: "host-a", Port: 80, AckedBy: "bob", AckedAt: time.Now()})
	_ = s.Append(Acknowledgment{ID: "a2", Host: "host-b", Port: 22, AckedBy: "alice", AckedAt: time.Now()})
	_ = s.Append(Acknowledgment{ID: "a3", Host: "host-a", Port: 443, AckedBy: "bob", AckedAt: time.Now()})

	result, err := s.LoadByHost("host-a")
	if err != nil {
		t.Fatalf("LoadByHost: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 entries for host-a, got %d", len(result))
	}
}

func TestAcknowledgmentStore_Delete(t *testing.T) {
	s := tempAcknowledgmentStore(t)
	_ = s.Append(Acknowledgment{ID: "a1", Host: "host-a", Port: 80, AckedBy: "bob", AckedAt: time.Now()})
	_ = s.Append(Acknowledgment{ID: "a2", Host: "host-a", Port: 443, AckedBy: "alice", AckedAt: time.Now()})

	if err := s.Delete("host-a", 80); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	entries, _ := s.Load()
	if len(entries) != 1 {
		t.Errorf("expected 1 entry after delete, got %d", len(entries))
	}
	if entries[0].Port != 443 {
		t.Errorf("expected remaining port 443, got %d", entries[0].Port)
	}
}

func TestAcknowledgmentStore_Append_InvalidDir(t *testing.T) {
	s := NewAcknowledgmentStore("/nonexistent/dir/acks.json")
	err := s.Append(Acknowledgment{ID: "x", Host: "h", Port: 1, AckedBy: "u", AckedAt: time.Now()})
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

func TestAcknowledgmentStore_MultipleEntries(t *testing.T) {
	s := tempAcknowledgmentStore(t)
	for i := 0; i < 5; i++ {
		_ = s.Append(Acknowledgment{
			ID:      string(rune('a' + i)),
			Host:    "10.0.0.1",
			Port:    8000 + i,
			AckedBy: "ops",
			AckedAt: time.Now(),
		})
	}
	entries, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 5 {
		t.Errorf("expected 5 entries, got %d", len(entries))
	}
}

func TestAcknowledgmentStore_WithExpiry(t *testing.T) {
	s := tempAcknowledgmentStore(t)
	expiry := time.Now().Add(24 * time.Hour).UTC().Truncate(time.Second)
	a := Acknowledgment{
		ID:        "ack-exp",
		Host:      "10.0.0.2",
		Port:      22,
		AckedBy:   "carol",
		Comment:   "temp ack",
		AckedAt:   time.Now().UTC().Truncate(time.Second),
		ExpiresAt: expiry,
	}
	_ = s.Append(a)
	entries, _ := s.Load()
	if entries[0].ExpiresAt.IsZero() {
		t.Error("expected non-zero ExpiresAt")
	}
	_ = os.Remove(s.path)
}
