package history

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

// PrintNotes writes notes for the given host to stdout.
// If host is empty, all notes are printed.
func PrintNotes(storePath, host string) error {
	s := NewNoteStore(storePath)
	notes, err := s.Load(host)
	if err != nil {
		return fmt.Errorf("loading notes: %w", err)
	}
	if len(notes) == 0 {
		fmt.Println("no notes found")
		return nil
	}
	return printNotesTable(os.Stdout, notes)
}

func printNotesTable(w io.Writer, notes []Note) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "HOST\tCREATED\tNOTE")
	for _, n := range notes {
		fmt.Fprintf(tw, "%s\t%s\t%s\n",
			n.Host,
			n.CreatedAt.Format("2006-01-02 15:04:05"),
			n.Text,
		)
	}
	return tw.Flush()
}
