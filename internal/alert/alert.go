package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Event represents a single alert event emitted on a port change.
type Event struct {
	Timestamp time.Time
	Host      string
	Level     Level
	Message   string
}

// Notifier writes alert events to an output destination.
type Notifier struct {
	Out io.Writer
}

// New returns a Notifier writing to stdout by default.
func New() *Notifier {
	return &Notifier{Out: os.Stdout}
}

// Notify emits alert events derived from a snapshot diff.
func (n *Notifier) Notify(host string, diff snapshot.Diff) {
	for _, p := range diff.Opened {
		n.emit(Event{
			Timestamp: time.Now(),
			Host:      host,
			Level:     LevelAlert,
			Message:   fmt.Sprintf("port %d newly OPEN", p),
		})
	}
	for _, p := range diff.Closed {
		n.emit(Event{
			Timestamp: time.Now(),
			Host:      host,
			Level:     LevelWarn,
			Message:   fmt.Sprintf("port %d newly CLOSED", p),
		})
	}
}

func (n *Notifier) emit(e Event) {
	fmt.Fprintf(n.Out, "[%s] %s %s — %s\n",
		e.Timestamp.Format(time.RFC3339),
		e.Level,
		e.Host,
		e.Message,
	)
}
