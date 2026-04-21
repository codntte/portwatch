package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/portwatch/internal/history"
)

func metricFile(dataDir string) string {
	return filepath.Join(dataDir, "metrics.jsonl")
}

func runMetric(args []string) {
	fs := flag.NewFlagSet("metric", flag.ExitOnError)
	host := fs.String("host", "", "filter by host")
	name := fs.String("name", "", "filter by metric name")
	dataDir := fs.String("data", ".", "data directory")
	append_ := fs.Bool("append", false, "append a metric entry")
	value := fs.Float64("value", 0, "metric value (used with -append)")
	unit := fs.String("unit", "", "metric unit (used with -append)")
	_ = fs.Parse(args)

	path := metricFile(*dataDir)

	if *append_ {
		if *host == "" || *name == "" {
			fmt.Fprintln(os.Stderr, "metric append requires -host and -name")
			os.Exit(1)
		}
		store := history.NewMetricStore(path)
		err := store.Append(history.MetricEntry{
			Host:  *host,
			Name:  *name,
			Value: *value,
			Unit:  *unit,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "error appending metric: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("metric recorded")
		return
	}

	if err := history.PrintMetrics(path, *host, *name); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
