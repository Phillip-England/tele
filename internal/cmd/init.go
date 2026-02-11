package cmd

import (
	"fmt"
	"os"

	"tele/internal/crypto"
	"tele/internal/store"
)

func RunInit() {
	exists, err := store.MasterExists()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error checking master config: %v\n", err)
		os.Exit(1)
	}
	if exists {
		fmt.Fprintln(os.Stderr, "Master password already configured. Delete master.json to reinitialize.")
		os.Exit(1)
	}

	fmt.Print("Enter master password: ")
	password, err := readPassword()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nError reading password: %v\n", err)
		os.Exit(1)
	}
	fmt.Println()

	if len(password) < 1 {
		fmt.Fprintln(os.Stderr, "Password cannot be empty.")
		os.Exit(1)
	}

	fmt.Print("Confirm master password: ")
	confirm, err := readPassword()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nError reading password: %v\n", err)
		os.Exit(1)
	}
	fmt.Println()

	if password != confirm {
		fmt.Fprintln(os.Stderr, "Passwords do not match.")
		os.Exit(1)
	}

	salt, err := crypto.GenerateSalt()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating salt: %v\n", err)
		os.Exit(1)
	}

	hash := crypto.HashPassword(password, salt)

	if err := store.WriteMaster(salt, hash); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing master config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Master password set successfully.")
}
