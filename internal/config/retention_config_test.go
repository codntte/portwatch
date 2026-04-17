package config

import (
	"testing"
	"time"
)

func TestRetentionConfig_MaxAge(t *testing.T) {
	r := RetentionConfig{MaxAgeDays: 7, MaxEntries: 500}
	if r.MaxAge() != 7*24*time.Hour {
		t.Errorf("expected 7 days, got %v", r.MaxAge())
	}
}

func TestRetentionConfig_MaxAge_Zero(t *testing.T) {
	r := RetentionConfig{MaxAgeDays: 0}
	if r.MaxAge() != 0 {
		t.Errorf("expected 0, got %v", r.MaxAge())
	}
}

func TestDefaultRetention(t *testing.T) {
	r := defaultRetention()
	if r.MaxAgeDays != 7 {
		t.Errorf("expected 7 days, got %d", r.MaxAgeDays)
	}
	if r.MaxEntries != 1000 {
		t.Errorf("expected 1000 entries, got %d", r.MaxEntries)
	}
}

func TestLoad_RetentionDefaults(t *testing.T) {
	path := writeTempConfig(t, `
hosts:
  - name: local
    address: localhost
    port_range: "80-81"
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Retention.MaxAgeDays != 7 {
		t.Errorf("expected default MaxAgeDays=7, got %d", cfg.Retention.MaxAgeDays)
	}
	if cfg.Retention.MaxEntries != 1000 {
		t.Errorf("expected default MaxEntries=1000, got %d", cfg.Retention.MaxEntries)
	}
}
