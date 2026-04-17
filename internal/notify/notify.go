package notify

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// Channel represents a notification output channel.
type Channel interface {
	Send(subject, body string) error
}

// StdoutChannel writes notifications to stdout (or any writer).
type StdoutChannel struct {
	Writer io.Writer
}

func NewStdoutChannel() *StdoutChannel {
	return &StdoutChannel{Writer: os.Stdout}
}

func (s *StdoutChannel) Send(subject, body string) error {
	_, err := fmt.Fprintf(s.Writer, "[%s] %s\n%s\n", time.Now().Format(time.RFC3339), subject, body)
	return err
}

// Dispatcher sends alert diffs through one or more channels.
type Dispatcher struct {
	channels []Channel
}

func NewDispatcher(channels ...Channel) *Dispatcher {
	return &Dispatcher{channels: channels}
}

// Dispatch formats and sends a diff via all registered channels.
func (d *Dispatcher) Dispatch(host string, diff alert.Diff) error {
	if len(diff.Opened) == 0 && len(diff.Closed) == 0 {
		return nil
	}

	subject := fmt.Sprintf("Port change detected on %s", host)
	body := formatBody(diff)

	for _, ch := range d.channels {
		if err := ch.Send(subject, body); err != nil {
			return fmt.Errorf("notify: channel send failed: %w", err)
		}
	}
	return nil
}

func formatBody(diff alert.Diff) string {
	body := ""
	for _, p := range diff.Opened {
		body += fmt.Sprintf("  [OPENED] port %d\n", p)
	}
	for _, p := range diff.Closed {
		body += fmt.Sprintf("  [CLOSED] port %d\n", p)
	}
	return body
}
