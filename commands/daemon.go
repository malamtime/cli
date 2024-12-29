package commands

import (
	"archive/zip"
	"embed"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gookit/color"
	"github.com/urfave/cli/v2"
)

//go:embed sys-desc/*
var sysDescFS embed.FS

var DaemonCommand *cli.Command = &cli.Command{
	Name:   "daemon:install",
	Usage:  "Install shelltime daemon service",
	Action: commandDaemonInstall,
}

func commandDaemonInstall(c *cli.Context) error {
	// Check if running as root
	if os.Geteuid() != 0 {
		return fmt.Errorf("this command must be run as root (sudo shelltime daemon:install)")
	}
	color.Yellow.Println("🔍 Detecting system architecture...")

	err := downloadDaemon()
	if err != nil {
		return err
	}

	switch runtime.GOOS {
	case "linux":
		return installLinuxDaemon()
	case "darwin":
		return installDarwinDaemon()
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

func downloadDaemon() error {
	arch := runtime.GOARCH
	goos := runtime.GOOS

	if goos == "windows" {
		return fmt.Errorf("windows is not supported for daemon installation")
	}

	// Determine archive extension
	ext := ".zip"
	if goos == "linux" {
		ext = ".tar.gz"
	}

	// Construct download URL
	downloadURL := fmt.Sprintf(
		"https://github.com/malamtime/cli/releases/latest/download/cli_%s_%s%s",
		strings.ToTitle(goos),
		arch,
		ext,
	)

	// Create temp dir for download
	tmpDir, err := os.MkdirTemp("", "shelltime-daemon")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Download archive
	color.Yellow.Printf("📥 Downloading daemon from %s...\n", downloadURL)
	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download daemon: %w", err)
	}
	defer resp.Body.Close()

	archivePath := filepath.Join(tmpDir, "archive"+ext)
	out, err := os.Create(archivePath)
	if err != nil {
		return fmt.Errorf("failed to create archive file: %w", err)
	}
	defer out.Close()

	if _, err = io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("failed to save archive: %w", err)
	}

	// Extract archive
	color.Yellow.Println("📦 Extracting daemon binary...")
	if goos == "linux" {
		cmd := exec.Command("tar", "xzf", archivePath, "-C", tmpDir)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to extract archive: %w", err)
		}
	} else {
		reader, err := zip.OpenReader(archivePath)
		if err != nil {
			return fmt.Errorf("failed to open archive: %w", err)
		}
		defer reader.Close()

		for _, file := range reader.File {
			if file.Name == "shelltime-daemon" {
				rc, err := file.Open()
				if err != nil {
					return fmt.Errorf("failed to open daemon binary from archive: %w", err)
				}
				defer rc.Close()

				dest := filepath.Join(tmpDir, "shelltime-daemon")
				out, err := os.Create(dest)
				if err != nil {
					return fmt.Errorf("failed to create daemon binary: %w", err)
				}
				defer out.Close()

				if _, err = io.Copy(out, rc); err != nil {
					return fmt.Errorf("failed to extract daemon binary: %w", err)
				}
				break
			}
		}
	}

	// Copy to final location
	binaryPath := "/usr/local/bin/shelltime-daemon"
	selfPath := filepath.Join(tmpDir, "shelltime-daemon")

	color.Yellow.Println("📦 Installing daemon binary...")
	if err := copyFile(selfPath, binaryPath); err != nil {
		return fmt.Errorf("failed to copy binary: %w", err)
	}
	if err := os.Chmod(binaryPath, 0755); err != nil {
		return fmt.Errorf("failed to set binary permissions: %w", err)
	}
	return nil
}

func installLinuxDaemon() error {
	serviceContent, err := sysDescFS.ReadFile("sys-desc/shelltime.service")
	if err != nil {
		return fmt.Errorf("failed to read service template: %w", err)
	}

	color.Yellow.Println("📝 Installing systemd service...")
	servicePath := "/etc/systemd/system/shelltime.service"
	if err := os.WriteFile(servicePath, serviceContent, 0644); err != nil {
		return fmt.Errorf("failed to write service file: %w", err)
	}

	color.Yellow.Println("🔄 Reloading systemd...")
	if err := exec.Command("systemctl", "daemon-reload").Run(); err != nil {
		return fmt.Errorf("failed to reload systemd: %w", err)
	}

	color.Yellow.Println("✨ Enabling service...")
	if err := exec.Command("systemctl", "enable", "shelltime").Run(); err != nil {
		return fmt.Errorf("failed to enable service: %w", err)
	}

	color.Yellow.Println("🚀 Starting service...")
	if err := exec.Command("systemctl", "start", "shelltime").Run(); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}

	color.Green.Println("✅ Daemon service has been installed and started successfully!")
	color.Green.Println("💡 Your commands will now be automatically synced to shelltime.xyz")
	return nil
}

func installDarwinDaemon() error {
	plistContent, err := sysDescFS.ReadFile("sys-desc/xyz.shelltime.daemon.plist")
	if err != nil {
		return fmt.Errorf("failed to read plist template: %w", err)
	}

	color.Yellow.Println("📝 Installing LaunchDaemon...")
	plistPath := "/Library/LaunchDaemons/xyz.shelltime.daemon.plist"
	if err := os.WriteFile(plistPath, plistContent, 0644); err != nil {
		return fmt.Errorf("failed to write plist file: %w", err)
	}

	color.Yellow.Println("🚀 Loading service...")
	if err := exec.Command("launchctl", "load", plistPath).Run(); err != nil {
		return fmt.Errorf("failed to load service: %w", err)
	}

	color.Green.Println("✅ Daemon service has been installed and started successfully!")
	color.Green.Println("💡 Your commands will now be automatically synced to shelltime.xyz")
	return nil
}

func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, input, 0644)
}
