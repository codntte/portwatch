package config

import "time"

const (
	// DefaultTimeout is used when a host timeout is not specified.
	DefaultTimeout = 2 * time.Second

	// DefaultInterval is the default scan interval.
	DefaultInterval = 5 * time.Minute

	// DefaultSnapshotDir is where snapshots are stored by default.
	DefaultSnapshotDir = ".portwatch"
)
