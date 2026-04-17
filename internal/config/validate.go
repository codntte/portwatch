package config

import (
	"errors"
	"fmt"
	"net"
)

// Validate checks the configuration for required fields and logical errors.
func Validate(cfg *Config) error {
	if len(cfg.Hosts) == 0 {
		return errors.New("config: at least one host must be specified")
	}
	for i, h := range cfg.Hosts {
		if h.Address == "" {
			return fmt.Errorf("config: host[%d] missing address", i)
		}
		if net.ParseIP(h.Address) == nil {
			// Allow hostnames — just ensure non-empty, resolved at scan time.
			if h.Address == "" {
				return fmt.Errorf("config: host[%d] invalid address %q", i, h.Address)
			}
		}
		if h.PortRange == "" {
			return fmt.Errorf("config: host[%d] missing port_range", i)
		}
		if h.Name == "" {
			cfg.Hosts[i].Name = h.Address
		}
		if h.Timeout == 0 {
			cfg.Hosts[i].Timeout = DefaultTimeout
		}
	}
	return nil
}
