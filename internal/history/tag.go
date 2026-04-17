package history

import (
	"encoding/json"
	"os"
	"time"
)

// Tag represents a named marker attached to a point in time.
type Tag struct {
	Name      string    `json:"name"`
	Note      string    `json:"note,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// TagStore manages tags stored in a JSON file.
type TagStore struct {
	path string
}

// NewTagStore returns a TagStore backed by the given file path.
func NewTagStore(path string) *TagStore {
	return &TagStore{path: path}
}

// Add appends a tag to the store.
func (s *TagStore) Add(name, note string) error {
	tags, err := s.Load()
	if err != nil {
		return err
	}
	tags = append(tags, Tag{Name: name, Note: note, CreatedAt: time.Now().UTC()})
	return s.save(tags)
}

// Load reads all tags from disk. Returns empty slice if file missing.
func (s *TagStore) Load() ([]Tag, error) {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return []Tag{}, nil
	}
	if err != nil {
		return nil, err
	}
	var tags []Tag
	if err := json.Unmarshal(data, &tags); err != nil {
		return nil, err
	}
	return tags, nil
}

// Delete removes all tags with the given name.
func (s *TagStore) Delete(name string) error {
	tags, err := s.Load()
	if err != nil {
		return err
	}
	filtered := tags[:0]
	for _, t := range tags {
		if t.Name != name {
			filtered = append(filtered, t)
		}
	}
	return s.save(filtered)
}

func (s *TagStore) save(tags []Tag) error {
	data, err := json.MarshalIndent(tags, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}
