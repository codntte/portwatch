package config

import (
	"os"
	"testing"
	"time"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestLoad_Valid(t *testing.T) {
	path := writeTempConfig(t, `
hosts:
  - name: localhost
    address: 127.0.0.1
    port_range: "22-80"
    timeout: 2s
interval: 10m
snapshot_dir: /tmp/snaps
alert:
  log_file: /tmp/portwatch.log
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Hosts) != 1 {
		t.Fatalf("expected 1 host, got %d", len(cfg.Hosts))
	}
	if cfg.Hosts[0].Address != "127.0.0.1" {
		t.Errorf("unexpected address: %s", cfg.Hosts[0].Address)
	}
	if cfg.Interval != 10*time.Minute {
		t.Errorf("unexpected interval: %v", cfg.Interval)
	}
	if cfg.Snapshot != "/tmp/snaps" {
		t.Errorf("unexpected snapshot_dir: %s", cfg.Snapshot)
	}
}

func TestLoad_Defaults(t *testing.T) {
	path := writeTempConfig(t, `hosts: []\n`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Interval != 5*time.Minute {
		t.Errorf("expected default interval 5m, got %v", cfg.Interval)
	}
	if cfg.Snapshot != ".portwatch" {
		t.Errorf("expected default snapshot_dir, got %s", cfg.Snapshot)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
