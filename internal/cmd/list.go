package cmd

import (
	"fmt"
	"os"

	"tele/internal/store"
)

func RunList() {
	names, err := store.ListDestinations()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing destinations: %v\n", err)
		os.Exit(1)
	}
	if len(names) == 0 {
		fmt.Println("No destinations saved.")
		return
	}
	for _, name := range names {
		host, port, user, _, _, _, err := store.ReadDestination(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", name, err)
			continue
		}
		fmt.Printf("  %s → %s@%s:%s\n", name, user, host, port)
	}
}
