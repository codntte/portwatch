package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/runner"
)

func main() {
	cfgPath := flag.String("config", "portwatch.yaml", "path to config file")
	once := flag.Bool("once", false, "run a single scan cycle and exit")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	if err := config.Validate(cfg); err != nil {
		log.Fatalf("config invalid: %v", err)
	}

	if err := os.MkdirAll(".portwatch", 0o755); err != nil {
		log.Fatalf("mkdir: %v", err)
	}

	channels := []notify.Channel{notify.NewStdoutChannel()}
	dispatcher := notify.NewDispatcher(channels)
	r := runner.New(cfg, dispatcher)

	if *once {
		if err := r.Run(); err != nil {
			log.Fatalf("run: %v", err)
		}
		return
	}

	interval := time.Duration(cfg.IntervalSeconds) * time.Second
	sched := runner.NewScheduler(r, interval)
	sched.Start()
}
