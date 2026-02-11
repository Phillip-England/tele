package cmd

import (
	"fmt"
	"os"

	"tele/internal/store"
)

func RunRm(name string) {
	if err := store.RemoveDestination(name); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Destination %q removed.\n", name)
}
