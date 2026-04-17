package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the full portwatch configuration.
type Config struct {
	Hosts    []HostConfig  `yaml:"hosts"`
	Alert    AlertConfig   `yaml:"alert"`
	Interval time.Duration `yaml:"interval"`
	Snapshot string        `yaml:"snapshot_dir"`
}

// HostConfig describes a single host to monitor.
type HostConfig struct {
	Name      string `yaml:"name"`
	Address   string `yaml:"address"`
	PortRange string `yaml:"port_range"`
	Timeout   time.Duration `yaml:"timeout"`
}

// AlertConfig holds alerting settings.
type AlertConfig struct {
	SlackWebhook string `yaml:"slack_webhook"`
	Email        string `yaml:"email"`
	LogFile      string `yaml:"log_file"`
}

// Load reads and parses a YAML config file.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	dec := yaml.NewDecoder(f)
	dec.KnownFields(true)
	if err := dec.Decode(&cfg); err != nil {
		return nil, err
	}

	if cfg.Snapshot == "" {
		cfg.Snapshot = ".portwatch"
	}
	if cfg.Interval == 0 {
		cfg.Interval = 5 * time.Minute
	}
	return &cfg, nil
}
