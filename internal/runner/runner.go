package runner

import (
	"log"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

// Runner orchestrates a single scan cycle for all configured hosts.
type Runner struct {
	cfg        *config.Config
	dispatcher *notify.Dispatcher
}

// New creates a Runner from the given config.
func New(cfg *config.Config, dispatcher *notify.Dispatcher) *Runner {
	return &Runner{cfg: cfg, dispatcher: dispatcher}
}

// Run executes one full scan cycle: scan → compare → alert.
func (r *Runner) Run() error {
	for _, host := range r.cfg.Hosts {
		if err := r.scanHost(host); err != nil {
			log.Printf("[portwatch] error scanning host %s: %v", host.Address, err)
		}
	}
	return nil
}

func (r *Runner) scanHost(host config.Host) error {
	opts := scanner.DefaultOptions()
	opts.Timeout = time.Duration(r.cfg.TimeoutSeconds) * time.Second

	ports, err := scanner.ScanPorts(host.Address, host.PortRange, opts)
	if err != nil {
		return err
	}

	snapshotPath := snapshotFile(host.Address)
	prev, _ := snapshot.Load(snapshotPath)

	curr := snapshot.New(host.Address, ports)
	if err := snapshot.Save(curr, snapshotPath); err != nil {
		return err
	}

	diff := snapshot.Compare(prev, curr)
	if diff.HasChanges() {
		body := alert.BuildDiff(diff)
		return r.dispatcher.Dispatch(host.Address, body)
	}
	return nil
}

func snapshotFile(address string) string {
	safe := ""
	for _, c := range address {
		if c == '.' || c == ':' || c == '/' {
			safe += "_"
		} else {
			safe += string(c)
		}
	}
	return ".portwatch/" + safe + ".json"
}
