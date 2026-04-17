package history

import (
	"fmt"
	"sort"
	"time"
)

// HostSummary holds aggregated change stats for a single host.
type HostSummary struct {
	Host        string
	TotalEvents int
	Opened      int
	Closed      int
	LastSeen    time.Time
}

// Summarize returns per-host summaries for entries matching the query.
func (s *Store) Summarize(q Query) ([]HostSummary, error) {
	entries, err := s.Query(q)
	if err != nil {
		return nil, fmt.Errorf("summarize: %w", err)
	}

	index := map[string]*HostSummary{}

	for _, e := range entries {
		hs, ok := index[e.Host]
		if !ok {
			hs = &HostSummary{Host: e.Host}
			index[e.Host] = hs
		}
		hs.TotalEvents++
		hs.Opened += len(e.Opened)
		hs.Closed += len(e.Closed)
		if e.Timestamp.After(hs.LastSeen) {
			hs.LastSeen = e.Timestamp
		}
	}

	result := make([]HostSummary, 0, len(index))
	for _, hs := range index {
		result = append(result, *hs)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Host < result[j].Host
	})
	return result, nil
}
