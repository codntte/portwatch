package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/user/portwatch/internal/history"
)

func cooldownFile(dir string) string {
	return dir + "/cooldown.json"
}

func runCooldown(args []string, dataDir string) {
	fs := flag.NewFlagSet("cooldown", flag.ExitOnError)
	host := fs.String("host", "", "Host to target")
	port := fs.Int("port", 0, "Port to target")
	duration := fs.Duration("duration", 30*time.Minute, "Cooldown duration (e.g. 30m, 2h)")
	delCmd := fs.Bool("delete", false, "Remove cooldown for host+port")
	listCmd := fs.Bool("list", false, "List all cooldown entries")
	_ = fs.Parse(args)

	store := history.NewCooldownStore(cooldownFile(dataDir))

	if *listCmd {
		entries, err := store.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error loading cooldowns: %v\n", err)
			os.Exit(1)
		}
		if len(entries) == 0 {
			fmt.Println("No cooldown entries found.")
			return
		}
		fmt.Printf("%-20s %-8s %-30s %s\n", "HOST", "PORT", "UNTIL", "ACTIVE")
		for _, e := range entries {
			active := "no"
			if time.Now().UTC().Before(e.Until) {
				active = "yes"
			}
			fmt.Printf("%-20s %-8s %-30s %s\n",
				e.Host,
				strconv.Itoa(e.Port),
				e.Until.Format(time.RFC3339),
				active,
			)
		}
		return
	}

	if *host == "" || *port == 0 {
		fmt.Fprintln(os.Stderr, "--host and --port are required")
		os.Exit(1)
	}

	if *delCmd {
		if err := store.Delete(*host, *port); err != nil {
			fmt.Fprintf(os.Stderr, "error deleting cooldown: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Cooldown removed for %s:%d\n", *host, *port)
		return
	}

	until := time.Now().UTC().Add(*duration)
	if err := store.Add(*host, *port, until); err != nil {
		fmt.Fprintf(os.Stderr, "error adding cooldown: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Cooldown set for %s:%d until %s\n", *host, *port, until.Format(time.RFC3339))
}
