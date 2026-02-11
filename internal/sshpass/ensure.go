package sshpass

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"tele/internal/config"
)

const sshpassURL = "https://sourceforge.net/projects/sshpass/files/sshpass/1.10/sshpass-1.10.tar.gz/download"

// BinDir returns the managed bin directory inside tele's config dir.
func BinDir() (string, error) {
	dir, err := config.Dir()
	if err != nil {
		return "", err
	}
	bin := filepath.Join(dir, "bin")
	if err := os.MkdirAll(bin, 0700); err != nil {
		return "", err
	}
	return bin, nil
}

// Ensure returns the path to a working sshpass binary.
// It checks PATH first, then tele's managed bin dir, and installs from source if needed.
func Ensure() (string, error) {
	// 1. Check if sshpass is already on PATH
	if p, err := exec.LookPath("sshpass"); err == nil {
		return p, nil
	}

	// 2. Check tele's managed bin dir
	binDir, err := BinDir()
	if err != nil {
		return "", fmt.Errorf("resolving bin dir: %w", err)
	}
	managed := filepath.Join(binDir, "sshpass")
	if _, err := os.Stat(managed); err == nil {
		return managed, nil
	}

	// 3. Install from source
	fmt.Println("sshpass not found — installing automatically...")
	if err := installFromSource(binDir); err != nil {
		return "", fmt.Errorf("auto-installing sshpass: %w", err)
	}

	if _, err := os.Stat(managed); err != nil {
		return "", fmt.Errorf("sshpass binary not found after install")
	}
	return managed, nil
}

func installFromSource(binDir string) error {
	// Verify a C compiler is available
	cc := findCC()
	if cc == "" {
		hint := "install Xcode Command Line Tools: xcode-select --install"
		if runtime.GOOS == "linux" {
			hint = "install gcc or clang via your package manager"
		}
		return fmt.Errorf("no C compiler found — %s", hint)
	}

	tmpDir, err := os.MkdirTemp("", "tele-sshpass-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	// Download tarball
	tarPath := filepath.Join(tmpDir, "sshpass.tar.gz")
	if err := download(sshpassURL, tarPath); err != nil {
		return fmt.Errorf("downloading sshpass: %w", err)
	}

	// Extract
	srcDir, err := extractTarGz(tarPath, tmpDir)
	if err != nil {
		return fmt.Errorf("extracting sshpass: %w", err)
	}

	// Configure and build
	configure := exec.Command("./configure", "--prefix="+tmpDir+"/out")
	configure.Dir = srcDir
	configure.Stdout = io.Discard
	configure.Stderr = io.Discard
	if err := configure.Run(); err != nil {
		return fmt.Errorf("configure failed: %w", err)
	}

	make := exec.Command("make")
	make.Dir = srcDir
	make.Stdout = io.Discard
	make.Stderr = io.Discard
	if err := make.Run(); err != nil {
		return fmt.Errorf("make failed: %w", err)
	}

	// Copy binary to managed bin dir
	src := filepath.Join(srcDir, "sshpass")
	dst := filepath.Join(binDir, "sshpass")
	if err := copyFile(src, dst); err != nil {
		return fmt.Errorf("copying binary: %w", err)
	}
	if err := os.Chmod(dst, 0755); err != nil {
		return err
	}

	fmt.Println("sshpass installed successfully.")
	return nil
}

func findCC() string {
	for _, name := range []string{"cc", "gcc", "clang"} {
		if p, err := exec.LookPath(name); err == nil {
			return p
		}
	}
	return ""
}

func download(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

// extractTarGz extracts a .tar.gz and returns the path to the top-level directory inside it.
func extractTarGz(tarPath, destDir string) (string, error) {
	f, err := os.Open(tarPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return "", err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	var topDir string
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		target := filepath.Join(destDir, hdr.Name)

		// Prevent path traversal
		if !strings.HasPrefix(filepath.Clean(target), filepath.Clean(destDir)+string(os.PathSeparator)) {
			continue
		}

		if topDir == "" {
			parts := strings.SplitN(hdr.Name, "/", 2)
			topDir = filepath.Join(destDir, parts[0])
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return "", err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return "", err
			}
			out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(hdr.Mode))
			if err != nil {
				return "", err
			}
			if _, err := io.Copy(out, tr); err != nil {
				out.Close()
				return "", err
			}
			out.Close()
		}
	}
	return topDir, nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
