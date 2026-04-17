package runner

import (
	"log"
	"time"
)

// Scheduler repeatedly calls Runner.Run on a fixed interval.
type Scheduler struct {
	runner   *Runner
	interval time.Duration
	stop     chan struct{}
}

// NewScheduler creates a Scheduler that triggers the runner at the given interval.
func NewScheduler(r *Runner, interval time.Duration) *Scheduler {
	return &Scheduler{
		runner:   r,
		interval: interval,
		stop:     make(chan struct{}),
	}
}

// Start begins the scheduling loop (blocking). Call Stop to exit.
// The runner is also invoked immediately on start before waiting for the
// first tick, so there is no delay on initial execution.
func (s *Scheduler) Start() {
	log.Printf("[portwatch] scheduler started, interval=%s", s.interval)
	if err := s.runner.Run(); err != nil {
		log.Printf("[portwatch] run error: %v", err)
	}
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := s.runner.Run(); err != nil {
				log.Printf("[portwatch] run error: %v", err)
			}
		case <-s.stop:
			log.Println("[portwatch] scheduler stopped")
			return
		}
	}
}

// Stop signals the scheduler to exit after the current run.
func (s *Scheduler) Stop() {
	close(s.stop)
}
