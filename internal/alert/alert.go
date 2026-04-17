package alert

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/user/portwatch/internal/snapshot"
)

// Diff holds the ports that were opened or closed between two snapshots.
type Diff struct {
	Opened []int
	Closed []int
}

// Notifier writes alerts for port changes.
type Notifier struct {
	w io.Writer
}

// New creates a Notifier writing to the given writer.
// If w is nil, os.Stdout is used.
func New(w io.Writer) *Notifier {
	if w == nil {
		w = os.Stdout
	}
	return &Notifier{w: w}
}

// Notify prints a human-readable diff to the notifier's writer.
func (n *Notifier) Notify(host string, diff snapshot.Diff) error {
	if len(diff.Opened) == 0 && len(diff.Closed) == 0 {
		return nil
	}

	opened := sortedKeys(diff.Opened)
	closed := sortedKeys(diff.Closed)

	for _, p := range opened {
		if _, err := fmt.Fprintf(n.w, "[%s] OPENED port %d\n", host, p); err != nil {
			return err
		}
	}
	for _, p := range closed {
		if _, err := fmt.Fprintf(n.w, "[%s] CLOSED port %d\n", host, p); err != nil {
			return err
		}
	}
	return nil
}

// BuildDiff converts a snapshot.Diff into an alert.Diff with sorted slices.
func BuildDiff(d snapshot.Diff) Diff {
	return Diff{
		Opened: sortedKeys(d.Opened),
		Closed: sortedKeys(d.Closed),
	}
}

func sortedKeys(m map[int]struct{}) []int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}
