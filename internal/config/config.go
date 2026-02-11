package config

import (
	"os"
	"path/filepath"
)

// Dir returns the tele config directory path, creating it if needed.
func Dir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(base, "tele")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}
	return dir, nil
}

// DestinationsDir returns the destinations subdirectory path, creating it if needed.
func DestinationsDir() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	dest := filepath.Join(dir, "destinations")
	if err := os.MkdirAll(dest, 0700); err != nil {
		return "", err
	}
	return dest, nil
}
