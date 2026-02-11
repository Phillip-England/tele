package main

import (
	"fmt"
	"os"

	"tele/internal/cmd"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		cmd.RunInit()
	case "add":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: tele add <name>")
			os.Exit(1)
		}
		cmd.RunAdd(os.Args[2])
	case "go":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: tele go <name>")
			os.Exit(1)
		}
		cmd.RunGo(os.Args[2])
	case "list":
		cmd.RunList()
	case "rm":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: tele rm <name>")
			os.Exit(1)
		}
		cmd.RunRm(os.Args[2])
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, `Usage: tele <command> [args]

Commands:
  init         Set up master password
  add <name>   Add a new SSH destination
  go <name>    SSH into a destination
  list         List all saved destinations
  rm <name>    Remove a destination`)
}
