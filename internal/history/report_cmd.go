package history

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
)

// PrintReport writes a human-readable report to stdout.
func PrintReport(diffFile, statsFile string) error {
	return printReportTo(os.Stdout, diffFile, statsFile)
}

func printReportTo(w io.Writer, diffFile, statsFile string) error {
	report, err := BuildReport(diffFile, statsFile)
	if err != nil {
		return fmt.Errorf("build report: %w", err)
	}

	if len(report.Entries) == 0 {
		fmt.Fprintln(w, "No report data available.")
		return nil
	}

	fmt.Fprintf(w, "Report generated at: %s\n\n", report.GeneratedAt.Format("2006-01-02 15:04:05 UTC"))

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "HOST\tOPEN PORTS\tCHANGES\tLAST SEEN")
	fmt.Fprintln(tw, "----\t----------\t-------\t---------")
	for _, e := range report.Entries {
		ports := formatReportPorts(e.OpenPorts)
		lastSeen := ""
		if !e.LastSeen.IsZero() {
			lastSeen = e.LastSeen.Format("2006-01-02 15:04:05")
		}
		fmt.Fprintf(tw, "%s\t%s\t%d\t%s\n", e.Host, ports, e.Changes, lastSeen)
	}
	return tw.Flush()
}

func formatReportPorts(ports []int) string {
	if len(ports) == 0 {
		return "-"
	}
	parts := make([]string, len(ports))
	for i, p := range ports {
		parts[i] = fmt.Sprintf("%d", p)
	}
	return strings.Join(parts, ",")
}
