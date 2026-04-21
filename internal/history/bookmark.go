package history

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

// Bookmark represents a named snapshot reference for a host at a point in time.
type Bookmark struct {
	Name      string    `json:"name"`
	Host      string    `json:"host"`
	CreatedAt time.Time `json:"created_at"`
	Note      string    `json:"note,omitempty"`
	Ports     []int     `json:"ports"`
}

// BookmarkStore manages named bookmarks persisted to a JSON file.
type BookmarkStore struct {
	path string
}

// NewBookmarkStore returns a BookmarkStore backed by the given file path.
func NewBookmarkStore(path string) *BookmarkStore {
	return &BookmarkStore{path: path}
}

// Add appends or replaces a bookmark with the given name.
func (s *BookmarkStore) Add(b Bookmark) error {
	marks, err := s.Load()
	if err != nil {
		return err
	}
	updated := false
	for i, m := range marks {
		if m.Name == b.Name && m.Host == b.Host {
			marks[i] = b
			updated = true
			break
		}
	}
	if !updated {
		marks = append(marks, b)
	}
	return s.save(marks)
}

// Load returns all bookmarks from the store.
func (s *BookmarkStore) Load() ([]Bookmark, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return []Bookmark{}, nil
	}
	if err != nil {
		return nil, err
	}
	var marks []Bookmark
	if err := json.Unmarshal(data, &marks); err != nil {
		return nil, err
	}
	return marks, nil
}

// Delete removes a bookmark by name and host.
func (s *BookmarkStore) Delete(name, host string) error {
	marks, err := s.Load()
	if err != nil {
		return err
	}
	filtered := marks[:0]
	for _, m := range marks {
		if m.Name != name || m.Host != host {
			filtered = append(filtered, m)
		}
	}
	return s.save(filtered)
}

func (s *BookmarkStore) save(marks []Bookmark) error {
	data, err := json.MarshalIndent(marks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}
