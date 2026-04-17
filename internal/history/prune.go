package history

import (
	"fmt"
	"os"
	"time"
)

// PruneOptions controls how old entries are removed.
type PruneOptions struct {
	// MaxAge removes entries older than this duration. Zero means no age limit.
	MaxAge time.Duration
	// MaxEntries keeps only the N most recent entries. Zero means no limit.
	MaxEntries int
}

// Prune removes entries from the host's history according to opts.
func (s *Store) Prune(host string, opts PruneOptions) error {
	entries, err := s.Load(host)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		return nil
	}

	now := time.Now().UTC()
	var filtered []Entry
	for _, e := range entries {
		if opts.MaxAge > 0 && now.Sub(e.Timestamp) > opts.MaxAge {
			continue
		}
		filtered = append(filtered, e)
	}

	if opts.MaxEntries > 0 && len(filtered) > opts.MaxEntries {
		filtered = filtered[len(filtered)-opts.MaxEntries:]
	}

	path := s.filePath(host)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("history: prune remove: %w", err)
	}
	for _, e := range filtered {
		if err := s.Append(e); err != nil {
			return err
		}
	}
	return nil
}
