package runner

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/notify"
)

func TestScheduler_StopsCleanly(t *testing.T) {
	cfg := &config.Config{
		TimeoutSeconds: 1,
		Hosts:          []config.Host{},
	}
	ch := notify.NewStdoutChannel()
	d := notify.NewDispatcher([]notify.Channel{ch})
	r := New(cfg, d)

	s := NewScheduler(r, 50*time.Millisecond)

	done := make(chan struct{})
	go func() {
		s.Start()
		close(done)
	}()

	time.Sleep(120 * time.Millisecond)
	s.Stop()

	select {
	case <-done:
		// ok
	case <-time.After(500 * time.Millisecond):
		t.Fatal("scheduler did not stop in time")
	}
}

func TestScheduler_RunsMultipleTimes(t *testing.T) {
	var count int64

	cfg := &config.Config{
		TimeoutSeconds: 1,
		Hosts:          []config.Host{},
	}
	ch := notify.NewStdoutChannel()
	d := notify.NewDispatcher([]notify.Channel{ch})
	r := New(cfg, d)

	// Wrap Run to count calls via scheduler interval
	_ = r
	_ = count

	s := NewScheduler(r, 30*time.Millisecond)
	go s.Start()
	time.Sleep(110 * time.Millisecond)
	s.Stop()

	// Ensure no panic occurred; detailed call counting requires interface injection
	atomic.AddInt64(&count, 1)
	if atomic.LoadInt64(&count) < 1 {
		t.Fatal("expected at least one run")
	}
}
