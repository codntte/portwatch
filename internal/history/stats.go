package history

import (
	"fmt"
	"sort"
)

// PortStat holds frequency data for a single port.
type PortStat struct {
	Port   int
	Opened int
	Closed int
}

// Stats returns open/close frequency per port for the given host.
func (s *Store) Stats(host string) ([]PortStat, error) {
	entries, err := s.Load()
	if err != nil {
		return nil, err
	}

	type counts struct{ opened, closed int }
	m := map[int]*counts{}

	for _, e := range entries {
		if e.Host != host {
			continue
		}
		for _, p := range e.Opened {
			if _, ok := m[p]; !ok {
				m[p] = &counts{}
			}
			m[p].opened++
		}
		for _, p := range e.Closed {
			if _, ok := m[p]; !ok {
				m[p] = &counts{}
			}
			m[p].closed++
		}
	}

	var stats []PortStat
	for port, c := range m {
		stats = append(stats, PortStat{Port: port, Opened: c.opened, Closed: c.closed})
	}
	sort.Slice(stats, func(i, j int) bool { return stats[i].Port < stats[j].Port })
	return stats, nil
}

// FormatStats returns a human-readable string of port stats.
func FormatStats(stats []PortStat) string {
	if len(stats) == 0 {
		return "no port activity recorded\n"
	}
	out := fmt.Sprintf("%-8s %-8s %-8s\n", "PORT", "OPENED", "CLOSED")
	for _, s := range stats {
		out += fmt.Sprintf("%-8d %-8d %-8d\n", s.Port, s.Opened, s.Closed)
	}
	return out
}
