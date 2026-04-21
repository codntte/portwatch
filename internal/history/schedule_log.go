package history

import (
	"encoding/json"
	"os"
	"time"
)

// ScheduleEntry records a single scan execution event.
type ScheduleEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Host      string    `json:"host"`
	Duration  string    `json:"duration"`
	PortsOpen int       `json:"ports_open"`
	Error     string    `json:"error,omitempty"`
}

// ScheduleLogStore manages a JSONL log of scan execution history.
type ScheduleLogStore struct {
	path string
}

// NewScheduleLogStore returns a new ScheduleLogStore backed by path.
func NewScheduleLogStore(path string) *ScheduleLogStore {
	return &ScheduleLogStore{path: path}
}

// Append writes a new ScheduleEntry to the log file.
func (s *ScheduleLogStore) Append(entry ScheduleEntry) error {
	f, err := os.OpenFile(s.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

// Load reads all ScheduleEntry records from the log file.
// Returns an empty slice if the file does not exist.
func (s *ScheduleLogStore) Load() ([]ScheduleEntry, error) {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return []ScheduleEntry{}, nil
	}
	if err != nil {
		return nil, err
	}

	var entries []ScheduleEntry
	for _, line := range splitLines(data) {
		var e ScheduleEntry
		if err := json.Unmarshal(line, &e); err != nil {
			continue
		}
		entries = append(entries, e)
	}
	return entries, nil
}
