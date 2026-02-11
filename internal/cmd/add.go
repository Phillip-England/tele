package cmd

import (
	"fmt"
	"os"

	"tele/internal/crypto"
	"tele/internal/store"
)

func RunAdd(name string) {
	exists, err := store.MasterExists()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if !exists {
		fmt.Fprintln(os.Stderr, "Not initialized. Run 'tele init' first.")
		os.Exit(1)
	}

	destExists, err := store.DestinationExists(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if destExists {
		fmt.Fprintf(os.Stderr, "Destination %q already exists.\n", name)
		os.Exit(1)
	}

	masterPass := verifyMasterPassword()

	host, err := promptLine("Host", "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if host == "" {
		fmt.Fprintln(os.Stderr, "Host cannot be empty.")
		os.Exit(1)
	}

	port, err := promptLine("Port", "22")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	user, err := promptLine("User", "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if user == "" {
		fmt.Fprintln(os.Stderr, "User cannot be empty.")
		os.Exit(1)
	}

	fmt.Print("Password: ")
	destPass, err := readPassword()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nError reading password: %v\n", err)
		os.Exit(1)
	}
	fmt.Println()

	salt, err := crypto.GenerateSalt()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating salt: %v\n", err)
		os.Exit(1)
	}

	key := crypto.DeriveKey(masterPass, salt)
	encPass, nonce, err := crypto.Encrypt([]byte(destPass), key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encrypting password: %v\n", err)
		os.Exit(1)
	}

	if err := store.WriteDestination(name, host, port, user, encPass, nonce, salt); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving destination: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Destination %q added.\n", name)
}
