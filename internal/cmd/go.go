package cmd

import (
	"fmt"
	"os"
	"syscall"

	"tele/internal/crypto"
	"tele/internal/sshpass"
	"tele/internal/store"
)

func RunGo(name string) {
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
	if !destExists {
		fmt.Fprintf(os.Stderr, "Destination %q not found.\n", name)
		os.Exit(1)
	}

	masterPass := verifyMasterPassword()

	host, port, user, encPass, nonce, salt, err := store.ReadDestination(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading destination: %v\n", err)
		os.Exit(1)
	}

	key := crypto.DeriveKey(masterPass, salt)
	destPass, err := crypto.Decrypt(encPass, nonce, key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decrypting password: %v\n", err)
		os.Exit(1)
	}

	sshpassPath, err := sshpass.Ensure()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	args := []string{
		"sshpass", "-p", string(destPass),
		"ssh",
		"-o", "StrictHostKeyChecking=no",
		"-p", port,
		fmt.Sprintf("%s@%s", user, host),
	}

	if err := syscall.Exec(sshpassPath, args, os.Environ()); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing ssh: %v\n", err)
		os.Exit(1)
	}
}
