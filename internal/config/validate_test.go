package config

import (
	"testing"
	"time"
)

func validConfig() *Config {
	return &Config{
		Hosts: []HostConfig{
			{Name: "local", Address: "127.0.0.1", PortRange: "1-1024", Timeout: 2 * time.Second},
		},
		Interval: 5 * time.Minute,
		Snapshot: ".portwatch",
	}
}

func TestValidate_Valid(t *testing.T) {
	if err := Validate(validConfig()); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_NoHosts(t *testing.T) {
	cfg := validConfig()
	cfg.Hosts = nil
	if err := Validate(cfg); err == nil {
		t.Error("expected error for empty hosts")
	}
}

func TestValidate_MissingAddress(t *testing.T) {
	cfg := validConfig()
	cfg.Hosts[0].Address = ""
	if err := Validate(cfg); err == nil {
		t.Error("expected error for missing address")
	}
}

func TestValidate_MissingPortRange(t *testing.T) {
	cfg := validConfig()
	cfg.Hosts[0].PortRange = ""
	if err := Validate(cfg); err == nil {
		t.Error("expected error for missing port_range")
	}
}

func TestValidate_DefaultsApplied(t *testing.T) {
	cfg := validConfig()
	cfg.Hosts[0].Name = ""
	cfg.Hosts[0].Timeout = 0
	if err := Validate(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Hosts[0].Name != cfg.Hosts[0].Address {
		t.Errorf("expected name to default to address")
	}
	if cfg.Hosts[0].Timeout != DefaultTimeout {
		t.Errorf("expected default timeout, got %v", cfg.Hosts[0].Timeout)
	}
}
