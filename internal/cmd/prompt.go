package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"

	"tele/internal/crypto"
	"tele/internal/store"
)

// readPassword reads a password from the terminal with echo disabled.
func readPassword() (string, error) {
	pw, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}
	return string(pw), nil
}

// readLine reads a line of text from stdin.
func readLine() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

// promptLine prints a prompt and reads a line. If the input is empty, returns defaultVal.
func promptLine(prompt, defaultVal string) (string, error) {
	if defaultVal != "" {
		fmt.Printf("%s [%s]: ", prompt, defaultVal)
	} else {
		fmt.Printf("%s: ", prompt)
	}
	val, err := readLine()
	if err != nil {
		return "", err
	}
	if val == "" {
		return defaultVal, nil
	}
	return val, nil
}

// verifyMasterPassword prompts for the master password and verifies it.
// Returns the password string on success, or exits on failure.
func verifyMasterPassword() string {
	fmt.Print("Enter master password: ")
	password, err := readPassword()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nError reading password: %v\n", err)
		os.Exit(1)
	}
	fmt.Println()

	salt, hash, err := store.ReadMaster()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading master config: %v\n", err)
		os.Exit(1)
	}

	if !crypto.VerifyPassword(password, salt, hash) {
		fmt.Fprintln(os.Stderr, "Incorrect master password.")
		os.Exit(1)
	}

	return password
}
