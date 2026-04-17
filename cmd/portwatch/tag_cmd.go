package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/portwatch/internal/history"
)

func tagFile(dataDir string) string {
	return filepath.Join(dataDir, "tags.json")
}

func runTag(args []string) {
	fs := flag.NewFlagSet("tag", flag.ExitOnError)
	add := fs.String("add", "", "add a tag with this name")
	del := fs.String("delete", "", "delete tags with this name")
	note := fs.String("note", "", "optional note for the tag")
	list := fs.Bool("list", false, "list all tags")
	dataDir := fs.String("data", ".", "directory for tag storage")
	_ = fs.Parse(args)

	store := history.NewTagStore(tagFile(*dataDir))

	switch {
	case *add != "":
		if err := store.Add(*add, *note); err != nil {
			fmt.Fprintf(os.Stderr, "error adding tag: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("tag %q added\n", *add)
	case *del != "":
		if err := store.Delete(*del); err != nil {
			fmt.Fprintf(os.Stderr, "error deleting tag: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("tag %q deleted\n", *del)
	case *list:
		if err := history.PrintTags(store, os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "error listing tags: %v\n", err)
			os.Exit(1)
		}
	default:
		fs.Usage()
	}
}
