// Package history records and manages per-host scan change history.
//
// Each host's history is stored as a newline-delimited JSON file (JSONL)
// under a configurable directory. Entries capture opened and closed ports
// with a UTC timestamp.
//
// Usage:
//
//	store := history.NewStore("/var/lib/portwatch/history")
//
//	// Record a change event.
//	store.Append(history.Entry{
//		Timestamp: time.Now().UTC(),
//		Host:      "192.168.1.1",
//		Opened:    []int{80},
//		Closed:    []int{22},
//	})
//
//	// Read all recorded events for a host.
//	entries, err := store.Load("192.168.1.1")
//
//	// Remove entries older than 7 days, keeping at most 100.
//	store.Prune("192.168.1.1", history.PruneOptions{
//		MaxAge:     7 * 24 * time.Hour,
//		MaxEntries: 100,
//	})
package history
