package history

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// PrintTags writes all tags from the store to w in a table format.
func PrintTags(store *TagStore, w io.Writer) error {
	tags, err := store.Load()
	if err != nil {
		return fmt.Errorf("load tags: %w", err)
	}
	if len(tags) == 0 {
		fmt.Fprintln(w, "no tags found")
		return nil
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "NAME\tCREATED\tNOTE")
	for _, t := range tags {
		fmt.Fprintf(tw, "%s\t%s\t%s\n", t.Name, t.CreatedAt.Format("2006-01-02 15:04:05"), t.Note)
	}
	return tw.Flush()
}
