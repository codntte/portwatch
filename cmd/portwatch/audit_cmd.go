package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/user/portwatch/internal/history"
)

func auditFile(dataDir string) string {
	return dataDir + "/audit.jsonl"
}

func runAudit(args []string) {
	fs := flag.NewFlagSet("audit", flag.ExitOnError)
	dataDir := fs.String("data", ".", "directory for data files")
	actor := fs.String("actor", "", "filter entries by actor")
	addActor := fs.String("add-actor", "", "actor name for new entry")
	addAction := fs.String("add-action", "add", "action for new entry (add|update|delete)")
	addTarget := fs.String("add-target", "", "target for new entry")
	addDetail := fs.String("add-detail", "", "optional detail for new entry")
	_ = fs.Parse(args)

	store := history.NewAuditStore(auditFile(*dataDir))

	if *addActor != "" {
		entry := history.AuditEntry{
			Actor:  *addActor,
			Action: history.AuditAction(*addAction),
			Target: *addTarget,
			Detail: *addDetail,
		}
		if err := store.Append(entry); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("audit entry recorded")
		return
	}

	var entries []history.AuditEntry
	var err error
	if *actor != "" {
		entries, err = store.LoadByActor(*actor)
	} else {
		entries, err = store.Load()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if len(entries) == 0 {
		fmt.Println("no audit entries found")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TIMESTAMP\tACTOR\tACTION\tTARGET\tDETAIL")
	for _, e := range entries {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			e.Timestamp.Format("2006-01-02 15:04:05"),
			e.Actor,
			e.Action,
			e.Target,
			e.Detail,
		)
	}
	w.Flush()
}
