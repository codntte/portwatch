package history

import (
	"encoding/json"
	"os"
	"time"
)

// WatchEvent represents a single port change event recorded during a watch cycle.
type WatchEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Host      string    `json:"host"`
	Opened    []int     `json:"opened,omitempty"`
	Closed    []int     `json:"closed,omitempty"`
}

// AppendEvent appends a WatchEvent to the history file at path.
func (s *Store) AppendEvent(e WatchEvent) error {
	f, err := os.OpenFile(s.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	line, err := json.Marshal(e)
	if err != nil {
		return err
	}
	_, err = f.Write(append(line, '\n'))
	return err
}

// LoadEvents reads all WatchEvents from the history file.
func (s *Store) LoadEvents() ([]WatchEvent, error) {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var events []WatchEvent
	for _, line := range splitLines(data) {
		var e WatchEvent
		if err := json.Unmarshal(line, &e); err != nil {
			continue
		}
		events = append(events, e)
	}
	return events, nil
}
