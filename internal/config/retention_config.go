package config

import "time"

// RetentionConfig holds history retention settings.
type RetentionConfig struct {
	MaxAgeDays int `yaml:"max_age_days"`
	MaxEntries int `yaml:"max_entries"`
}

// MaxAge converts MaxAgeDays to a time.Duration.
func (r RetentionConfig) MaxAge() time.Duration {
	return time.Duration(r.MaxAgeDays) * 24 * time.Hour
}

// defaultRetention returns default retention settings.
func defaultRetention() RetentionConfig {
	return RetentionConfig{
		MaxAgeDays: 7,
		MaxEntries: 1000,
	}
}
