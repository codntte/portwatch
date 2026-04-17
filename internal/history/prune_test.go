package history_test

import (
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/history"
)

func makeEntry(host string, ago time.Duration, opened []int) history.Entry {
	return history.Entry{
		Timestamp: time.Now().UTC().Add(-ago),
		Host:      host,
		Opened:    opened,
		Closed:    []int{},
	}
}

func TestPrune_MaxAge(t *testing.T) {
	store := history.NewStore(t.TempDir())
	_ = store.Append(makeEntry("h", 2*time.Hour, []int{80}))
	_ = store.Append(makeEntry("h", 30*time.Minute, []int{443}))
	_ = store.Append(makeEntry("h", 5*time.Minute, []int{8080}))

	if err := store.Prune("h", history.PruneOptions{MaxAge: time.Hour}); err != nil {
		t.Fatalf("Prune: %v", err)
	}
	entries, _ := store.Load("h")
	if len(entries) != 2 {
		t.Errorf("expected 2 entries after age prune, got %d", len(entries))
	}
}

func TestPrune_MaxEntries(t *testing.T) {
	store := history.NewStore(t.TempDir())
	for i := 0; i < 5; i++ {
		_ = store.Append(makeEntry("h2", time.Duration(i)*time.Minute, []int{i}))
	}

	if err := store.Prune("h2", history.PruneOptions{MaxEntries: 3}); err != nil {
		t.Fatalf("Prune: %v", err)
	}
	entries, _ := store.Load("h2")
	if len(entries) != 3 {
		t.Errorf("expected 3 entries after count prune, got %d", len(entries))
	}
}

func TestPrune_NoFile(t *testing.T) {
	store := history.NewStore(t.TempDir())
	if err := store.Prune("nobody", history.PruneOptions{MaxAge: time.Hour}); err != nil {
		t.Errorf("expected no error pruning missing host, got %v", err)
	}
}
