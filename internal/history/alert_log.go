package history

import (
	"encoding/json"
	"os"
	"time"
)

// AlertEntry records a single alert event for a host.
type AlertEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Host      string    `json:"host"`
	Opened    []int     `json:"opened,omitempty"`
	Closed    []int     `json:"closed,omitempty"`
}

// AlertLogStore persists alert entries to a JSONL file.
type AlertLogStore struct {
	path string
}

// NewAlertLogStore creates a new AlertLogStore backed by the given file path.
func NewAlertLogStore(path string) *AlertLogStore {
	return &AlertLogStore{path: path}
}

// Append writes a new AlertEntry to the log file.
func (s *AlertLogStore) Append(entry AlertEntry) error {
	f, err := os.OpenFile(s.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	line, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	_, err = f.Write(append(line, '\n'))
	return err
}

// Load reads all AlertEntry records from the log file.
// Returns an empty slice if the file does not exist.
func (s *AlertLogStore) Load() ([]AlertEntry, error) {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return []AlertEntry{}, nil
	}
	if err != nil {
		return nil, err
	}

	var entries []AlertEntry
	for _, line := range splitLines(data) {
		var e AlertEntry
		if err := json.Unmarshal(line, &e); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}
