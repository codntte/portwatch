package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/user/portwatch/internal/history"
)

func policyFile(dir string) string {
	return dir + "/policies.json"
}

func runPolicy(args []string, dataDir string) error {
	fs := flag.NewFlagSet("policy", flag.ContinueOnError)
	add := fs.String("add", "", "add or replace a policy by name")
	del := fs.String("delete", "", "delete a policy by name")
	host := fs.String("host", "", "host address for the policy")
	portRange := fs.String("ports", "", "port range (e.g. 1-1024)")
	interval := fs.Duration("interval", 0, "scan interval (e.g. 5m, 1h)")
	enabled := fs.Bool("enabled", true, "whether the policy is active")

	if err := fs.Parse(args); err != nil {
		return err
	}

	path := policyFile(dataDir)
	store := history.NewPolicyStore(path)

	switch {
	case *add != "":
		if *host == "" || *portRange == "" {
			return fmt.Errorf("policy: --host and --ports are required with --add")
		}
		entry := history.PolicyEntry{
			Name:      *add,
			Host:      *host,
			PortRange: *portRange,
			Interval:  *interval,
			Enabled:   *enabled,
		}
		if err := store.Add(entry); err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "policy %q saved\n", *add)

	case *del != "":
		if err := store.Delete(*del); err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "policy %q deleted\n", *del)

	default:
		if err := history.PrintPolicies(path); err != nil {
			return err
		}
	}

	_ = time.Second // keep import used via PolicyEntry.Interval
	return nil
}
