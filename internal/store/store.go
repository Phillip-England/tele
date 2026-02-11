package store

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"tele/internal/config"
)

// MasterConfig represents the master.json file on disk.
type MasterConfig struct {
	Salt         string `json:"salt"`
	PasswordHash string `json:"password_hash"`
}

// Destination represents a destination JSON file on disk.
type Destination struct {
	Host              string `json:"host"`
	Port              string `json:"port"`
	User              string `json:"user"`
	EncryptedPassword string `json:"encrypted_password"`
	Nonce             string `json:"nonce"`
	Salt              string `json:"salt"`
}

// MasterExists checks if master.json exists.
func MasterExists() (bool, error) {
	dir, err := config.Dir()
	if err != nil {
		return false, err
	}
	_, err = os.Stat(filepath.Join(dir, "master.json"))
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}

// WriteMaster writes the master config to disk.
func WriteMaster(salt, passwordHash []byte) error {
	dir, err := config.Dir()
	if err != nil {
		return err
	}
	mc := MasterConfig{
		Salt:         hex.EncodeToString(salt),
		PasswordHash: hex.EncodeToString(passwordHash),
	}
	data, err := json.MarshalIndent(mc, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "master.json"), data, 0600)
}

// ReadMaster reads the master config from disk, returning salt and passwordHash as bytes.
func ReadMaster() (salt, passwordHash []byte, err error) {
	dir, err := config.Dir()
	if err != nil {
		return nil, nil, err
	}
	data, err := os.ReadFile(filepath.Join(dir, "master.json"))
	if err != nil {
		return nil, nil, err
	}
	var mc MasterConfig
	if err := json.Unmarshal(data, &mc); err != nil {
		return nil, nil, err
	}
	salt, err = hex.DecodeString(mc.Salt)
	if err != nil {
		return nil, nil, fmt.Errorf("decoding salt: %w", err)
	}
	passwordHash, err = hex.DecodeString(mc.PasswordHash)
	if err != nil {
		return nil, nil, fmt.Errorf("decoding password hash: %w", err)
	}
	return salt, passwordHash, nil
}

// WriteDestination saves a destination to disk.
func WriteDestination(name string, host, port, user string, encPass, nonce, salt []byte) error {
	dir, err := config.DestinationsDir()
	if err != nil {
		return err
	}
	d := Destination{
		Host:              host,
		Port:              port,
		User:              user,
		EncryptedPassword: hex.EncodeToString(encPass),
		Nonce:             hex.EncodeToString(nonce),
		Salt:              hex.EncodeToString(salt),
	}
	data, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, name+".json"), data, 0600)
}

// ReadDestination reads a destination from disk.
func ReadDestination(name string) (host, port, user string, encPass, nonce, salt []byte, err error) {
	dir, err := config.DestinationsDir()
	if err != nil {
		return "", "", "", nil, nil, nil, err
	}
	data, err := os.ReadFile(filepath.Join(dir, name+".json"))
	if err != nil {
		return "", "", "", nil, nil, nil, err
	}
	var d Destination
	if err := json.Unmarshal(data, &d); err != nil {
		return "", "", "", nil, nil, nil, err
	}
	encPass, err = hex.DecodeString(d.EncryptedPassword)
	if err != nil {
		return "", "", "", nil, nil, nil, fmt.Errorf("decoding encrypted password: %w", err)
	}
	nonce, err = hex.DecodeString(d.Nonce)
	if err != nil {
		return "", "", "", nil, nil, nil, fmt.Errorf("decoding nonce: %w", err)
	}
	salt, err = hex.DecodeString(d.Salt)
	if err != nil {
		return "", "", "", nil, nil, nil, fmt.Errorf("decoding salt: %w", err)
	}
	return d.Host, d.Port, d.User, encPass, nonce, salt, nil
}

// ListDestinations returns the names of all saved destinations.
func ListDestinations() ([]string, error) {
	dir, err := config.DestinationsDir()
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".json") {
			names = append(names, strings.TrimSuffix(e.Name(), ".json"))
		}
	}
	return names, nil
}

// RemoveDestination deletes a destination file.
func RemoveDestination(name string) error {
	dir, err := config.DestinationsDir()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, name+".json")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("destination %q not found", name)
	}
	return os.Remove(path)
}

// DestinationExists checks if a destination file exists.
func DestinationExists(name string) (bool, error) {
	dir, err := config.DestinationsDir()
	if err != nil {
		return false, err
	}
	_, err = os.Stat(filepath.Join(dir, name+".json"))
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}
