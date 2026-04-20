package history

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

// Note represents a free-form user note attached to a host at a point in time.
type Note struct {
	Host      string    `json:"host"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

// NoteStore manages persistent notes stored as newline-delimited JSON.
type NoteStore struct {
	path string
}

// NewNoteStore returns a NoteStore backed by the given file path.
func NewNoteStore(path string) *NoteStore {
	return &NoteStore{path: path}
}

// Add appends a new note for the given host.
func (s *NoteStore) Add(host, text string) error {
	if host == "" {
		return errors.New("host must not be empty")
	}
	if text == "" {
		return errors.New("note text must not be empty")
	}
	n := Note{Host: host, Text: text, CreatedAt: time.Now().UTC()}
	b, err := json.Marshal(n)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(s.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(append(b, '\n'))
	return err
}

// Load returns all notes, optionally filtered by host (empty string = all).
func (s *NoteStore) Load(host string) ([]Note, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	var notes []Note
	for _, line := range splitLines(string(data)) {
		var n Note
		if err := json.Unmarshal([]byte(line), &n); err != nil {
			continue
		}
		if host == "" || n.Host == host {
			notes = append(notes, n)
		}
	}
	return notes, nil
}

// Delete removes all notes for the given host.
func (s *NoteStore) Delete(host string) error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	var kept []string
	for _, line := range splitLines(string(data)) {
		var n Note
		if err := json.Unmarshal([]byte(line), &n); err != nil || n.Host == host {
			continue
		}
		kept = append(kept, line)
	}
	var out string
	for _, l := range kept {
		out += l + "\n"
	}
	return os.WriteFile(s.path, []byte(out), 0o644)
}
