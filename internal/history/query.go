package history

import (
	"time"
)

// QueryOptions defines filters for querying history entries.
type QueryOptions struct {
	Host  string
	Since time.Time
	Until time.Time
	Limit int
}

// Query returns entries from the store matching the given options.
func (s *Store) Query(opts QueryOptions) ([]Entry, error) {
	entries, err := s.Load()
	if err != nil {
		return nil, err
	}

	var result []Entry
	for _, e := range entries {
		if opts.Host != "" && e.Host != opts.Host {
			continue
		}
		if !opts.Since.IsZero() && e.Timestamp.Before(opts.Since) {
			continue
		}
		if !opts.Until.IsZero() && e.Timestamp.After(opts.Until) {
			continue
		}
		result = append(result, e)
		if opts.Limit > 0 && len(result) >= opts.Limit {
			break
		}
	}
	return result, nil
}
