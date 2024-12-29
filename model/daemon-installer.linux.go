package model

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gookit/color"
)

//go:embed sys-desc/shelltime.service
var daemonLinuxServiceDesc []byte

// LinuxDaemonInstaller implements DaemonInstaller for Linux systems
type LinuxDaemonInstaller struct {
	baseFolder string
}

func NewLinuxDaemonInstaller(baseFolder string) *LinuxDaemonInstaller {
	return &LinuxDaemonInstaller{baseFolder: baseFolder}
}

func (l *LinuxDaemonInstaller) CheckAndStopExistingService() error {
	color.Yellow.Println("🔍 Checking if service is running...")
	cmd := exec.Command("systemctl", "is-active", "shelltime")
	if err := cmd.Run(); err == nil {
		color.Yellow.Println("🛑 Stopping existing service...")
		if err := exec.Command("systemctl", "stop", "shelltime").Run(); err != nil {
			return fmt.Errorf("failed to stop existing service: %w", err)
		}
	}
	return nil
}

func (l *LinuxDaemonInstaller) InstallService() error {
	daemonPath := filepath.Join(l.baseFolder, "daemon")
	// Create daemon directory if not exists
	if err := os.MkdirAll(daemonPath, 0755); err != nil {
		return fmt.Errorf("failed to create daemon directory: %w", err)
	}

	servicePath := filepath.Join(daemonPath, "shelltime.service")
	if _, err := os.Stat(servicePath); err == nil {
		if err := os.Remove(servicePath); err != nil {
			return fmt.Errorf("failed to remove existing service file: %w", err)
		}
	}

	if err := os.WriteFile(servicePath, daemonLinuxServiceDesc, 0644); err != nil {
		return fmt.Errorf("failed to write service file: %w", err)
	}

	return nil
}

func (l *LinuxDaemonInstaller) RegisterService() error {
	servicePath := "/etc/systemd/system/shelltime.service"
	if _, err := os.Stat(servicePath); err != nil {
		sourceFile := filepath.Join(l.baseFolder, "daemon/shelltime.service")
		if err := os.Symlink(sourceFile, servicePath); err != nil {
			return fmt.Errorf("failed to create service symlink: %w", err)
		}
	}
	return nil
}

func (l *LinuxDaemonInstaller) StartService() error {
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
	return nil
}

func (l *LinuxDaemonInstaller) UnregisterService() error {
	color.Yellow.Println("🛑 Stopping and disabling service if running...")
	// Try to stop and disable the service
	_ = exec.Command("systemctl", "stop", "shelltime").Run()
	_ = exec.Command("systemctl", "disable", "shelltime").Run()

	color.Yellow.Println("🗑 Removing service files...")
	// Remove symlink from systemd
	if err := os.Remove("/etc/systemd/system/shelltime.service"); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove systemd service symlink: %w", err)
	}

	color.Yellow.Println("🔄 Reloading systemd...")
	if err := exec.Command("systemctl", "daemon-reload").Run(); err != nil {
		return fmt.Errorf("failed to reload systemd: %w", err)
	}

	color.Green.Println("✅ Service unregistered successfully")
	return nil
}
