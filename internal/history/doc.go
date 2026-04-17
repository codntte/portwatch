// Package history provides persistent storage for port change events.
//
// It supports appending scan diff entries, loading historical records,
// pruning old entries by age or count, exporting to CSV, and querying
// entries by host, time range, or limit.
//
// Typical usage:
//
//	store := history.NewStore("/var/lib/portwatch/history.json")
//
//	// Append a new change event
//	_ = store.Append(history.Entry{
//		Host:      "192.168.1.1",
//		Timestamp: time.Now(),
//		Opened:    []int{80, 443},
//		Closed:    []int{22},
//	})
//
//	// Query recent changes for a specific host
//	entries, _ := store.Query(history.QueryOptions{
//		Host:  "192.168.1.1",
//		Since: time.Now().Add(-24 * time.Hour),
//	})
//
//	// Prune entries older than 30 days
//	_ = store.Prune(history.PruneOptions{MaxAge: 30 * 24 * time.Hour})
//
//	// Export to CSV
//	_ = store.ExportCSV("/tmp/history.csv")
package history
