package commands

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/gookit/color"
	"github.com/urfave/cli/v2"
)

//go:embed sys-desc/*
var sysDescFS embed.FS

var DaemonInstallCommand *cli.Command = &cli.Command{
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
	goos := runtime.GOOS

	if goos == "windows" {
		return fmt.Errorf("windows is not supported for daemon installation")
	}

	if goos == "linux" {
		color.Yellow.Println("🔍 Checking if service is running...")
		cmd := exec.Command("systemctl", "is-active", "shelltime")
		if err := cmd.Run(); err == nil {
			color.Yellow.Println("🛑 Stopping existing service...")
			if err := exec.Command("systemctl", "stop", "shelltime").Run(); err != nil {
				return fmt.Errorf("failed to stop existing service: %w", err)
			}
		}
	} else if goos == "darwin" {
		color.Yellow.Println("🔍 Checking if service is running...")
		cmd := exec.Command("launchctl", "list", "xyz.shelltime.daemon")
		if err := cmd.Run(); err == nil {
			color.Yellow.Println("🛑 Stopping existing service...")
			if err := exec.Command("launchctl", "unload", "/Library/LaunchDaemons/xyz.shelltime.daemon.plist").Run(); err != nil {
				return fmt.Errorf("failed to stop existing service: %w", err)
			}
		}
	}

	realBin, err := getRealBinPath()
	if err != nil {
		return err
	}

	// check backup file exist or not
	if _, err := os.Stat(filepath.Join(realBin, "bin/shelltime-daemon.bak")); err == nil {
		color.Yellow.Println("🔄 Found latest daemon file, restoring...")
		// try to remove old file
		_ = os.Remove(filepath.Join(realBin, "bin/shelltime-daemon"))
		// rename .bak to original
		if err := os.Rename(
			filepath.Join(realBin, "bin/shelltime-daemon.bak"),
			filepath.Join(realBin, "bin/shelltime-daemon"),
		); err != nil {
			return fmt.Errorf("failed to restore latest daemon: %w", err)
		}
	}

	// check shelltime-daemon
	if _, err := os.Stat(filepath.Join(realBin, "bin/shelltime-daemon")); err != nil {
		color.Yellow.Println("⚠️ shelltime-daemon not found, please reinstall the CLI first:")
		color.Yellow.Println("curl -sSL https://raw.githubusercontent.com/malamtime/installation/master/install.bash | bash")
		return nil
	}

	// Copy to final location
	binaryPath := "/usr/local/bin/shelltime-daemon"

	if _, err := os.Stat(binaryPath); err != nil {
		color.Yellow.Println("🔍 Creating daemon symlink...")
		if err := os.Symlink(filepath.Join(realBin, "bin/shelltime-daemon"), binaryPath); err != nil {
			return fmt.Errorf("failed to create daemon symlink: %w", err)
		}
	}

	if err := installDaemonDescriptionFilesLocally(realBin); err != nil {
		return err
	}

	if err := registerDaemonService(realBin); err != nil {
		return err
	}

	return startDaemonService()
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

// getRealBinPath will return the first matched `~/.shelltime/` folder
func getRealBinPath() (string, error) {
	homeAbsolutePrefix := ""
	var scanPaths []string
	if runtime.GOOS == "linux" {
		homeAbsolutePrefix = "/home"
	} else if runtime.GOOS == "darwin" {
		homeAbsolutePrefix = "/Users"
	}
	scanPaths = append(scanPaths, homeAbsolutePrefix)

	// Scan paths for .shelltime/bin folder
	foundUser := ""
	for _, basePath := range scanPaths {
		entries, err := os.ReadDir(basePath)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				shelltimePath := filepath.Join(basePath, entry.Name(), ".shelltime", "bin")
				if _, err := os.Stat(shelltimePath); err == nil {
					foundUser = entry.Name()
					break
				}
			}
		}
		if foundUser != "" {
			break
		}
	}

	if foundUser == "" && runtime.GOOS == "linux" {
		shelltimePath := filepath.Join("/root", ".shelltime", "bin")
		if _, err := os.Stat(shelltimePath); err == nil {
			foundUser = "root"
		}
	}

	if foundUser == "" {
		return "", fmt.Errorf("could not find any user with ~/.shelltime/bin directory")
	}

	if foundUser == "root" && runtime.GOOS == "linux" {
		return filepath.Join("/root", ".shelltime"), nil
	}

	return filepath.Join(homeAbsolutePrefix, foundUser, ".shelltime"), nil
}

func installDaemonDescriptionFilesLocally(baseFolder string) error {
	daemonPath := filepath.Join(baseFolder, "daemon")

	// Create daemon directory if not exists
	if err := os.MkdirAll(daemonPath, 0755); err != nil {
		return fmt.Errorf("failed to create daemon directory: %w", err)
	}

	// Copy Linux service file
	serviceContent, err := sysDescFS.ReadFile("sys-desc/shelltime.service")
	if err != nil {
		return fmt.Errorf("failed to read service template: %w", err)
	}

	servicePath := filepath.Join(daemonPath, "shelltime.service")
	if _, err := os.Stat(servicePath); err == nil {
		if err := os.Remove(servicePath); err != nil {
			return fmt.Errorf("failed to remove existing service file: %w", err)
		}
	}

	if err := os.WriteFile(servicePath, serviceContent, 0644); err != nil {
		return fmt.Errorf("failed to write service file: %w", err)
	}

	// Copy macOS plist file
	plistContent, err := sysDescFS.ReadFile("sys-desc/xyz.shelltime.daemon.plist")
	if err != nil {
		return fmt.Errorf("failed to read plist template: %w", err)
	}

	plistPath := filepath.Join(daemonPath, "xyz.shelltime.daemon.plist")
	if _, err := os.Stat(plistPath); err == nil {
		if err := os.Remove(plistPath); err != nil {
			return fmt.Errorf("failed to remove existing plist file: %w", err)
		}
	}

	if err := os.WriteFile(plistPath, plistContent, 0644); err != nil {
		return fmt.Errorf("failed to write plist file: %w", err)
	}

	return nil
}

func registerDaemonService(baseFolder string) error {
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	switch runtime.GOOS {
	case "linux":
		servicePath := "/etc/systemd/system/shelltime.service"
		if _, err := os.Stat(servicePath); err != nil {
			sourceFile := filepath.Join(baseFolder, "daemon/shelltime.service")
			if err := os.Symlink(sourceFile, servicePath); err != nil {
				return fmt.Errorf("failed to create service symlink: %w", err)
			}
		}
	case "darwin":
		plistPath := "/Library/LaunchDaemons/xyz.shelltime.daemon.plist"
		if _, err := os.Stat(plistPath); err != nil {
			sourceFile := filepath.Join(baseFolder, "daemon/xyz.shelltime.daemon.plist")
			if err := os.Symlink(sourceFile, plistPath); err != nil {
				return fmt.Errorf("failed to create plist symlink: %w", err)
			}
		}
	}

	return nil
}
func startDaemonService() error {
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	switch runtime.GOOS {
	case "linux":
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

	case "darwin":
		color.Yellow.Println("🚀 Starting service...")
		if err := exec.Command("launchctl", "load", "/Library/LaunchDaemons/xyz.shelltime.daemon.plist").Run(); err != nil {
			return fmt.Errorf("failed to start service: %w", err)
		}

	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	color.Green.Println("✅ Daemon service has been started successfully!")
	return nil
}
