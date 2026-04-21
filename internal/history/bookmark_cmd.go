package history

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
)

// PrintBookmarks writes a formatted table of bookmarks to stdout.
func PrintBookmarks(storePath, filterHost string) error {
	return printBookmarksTo(os.Stdout, storePath, filterHost)
}

func printBookmarksTo(w io.Writer, storePath, filterHost string) error {
	s := NewBookmarkStore(storePath)
	marks, err := s.Load()
	if err != nil {
		return fmt.Errorf("load bookmarks: %w", err)
	}

	if filterHost != "" {
		filtered := marks[:0]
		for _, m := range marks {
			if m.Host == filterHost {
				filtered = append(filtered, m)
			}
		}
		marks = filtered
	}

	if len(marks) == 0 {
		fmt.Fprintln(w, "no bookmarks found")
		return nil
	}

	sort.Slice(marks, func(i, j int) bool {
		if marks[i].Host != marks[j].Host {
			return marks[i].Host < marks[j].Host
		}
		return marks[i].Name < marks[j].Name
	})

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "NAME\tHOST\tPORTS\tCREATED\tNOTE")
	for _, m := range marks {
		ports := formatBookmarkPorts(m.Ports)
		created := m.CreatedAt.Format("2006-01-02 15:04")
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n",
			m.Name, m.Host, ports, created, m.Note)
	}
	return tw.Flush()
}

func formatBookmarkPorts(ports []int) string {
	if len(ports) == 0 {
		return "-"
	}
	parts := make([]string, len(ports))
	for i, p := range ports {
		parts[i] = fmt.Sprintf("%d", p)
	}
	return strings.Join(parts, ",")
}
