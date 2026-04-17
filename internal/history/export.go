package history

import (
	"encoding/csv"
	"fmt"
	"io"
	"time"
)

// ExportCSV writes all history entries for a given host to w in CSV format.
func (s *Store) ExportCSV(host string, w io.Writer) error {
	entries, err := s.Load(host)
	if err != nil {
		return fmt.Errorf("export csv: %w", err)
	}

	cw := csv.NewWriter(w)
	defer cw.Flush()

	if err := cw.Write([]string{"timestamp", "event", "port"}); err != nil {
		return fmt.Errorf("export csv header: %w", err)
	}

	for _, e := range entries {
		ts := e.Timestamp.Format(time.RFC3339)
		for _, p := range e.Opened {
			if err := cw.Write([]string{ts, "opened", fmt.Sprintf("%d", p)}); err != nil {
				return fmt.Errorf("export csv row: %w", err)
			}
		}
		for _, p := range e.Closed {
			if err := cw.Write([]string{ts, "closed", fmt.Sprintf("%d", p)}); err != nil {
				return fmt.Errorf("export csv row: %w", err)
			}
		}
	}

	return cw.Error()
}
