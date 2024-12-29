package model

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gookit/color"
)

//go:embed sys-desc/xyz.shelltime.daemon.plist
var daemonMacServiceDesc []byte

// MacDaemonInstaller implements DaemonInstaller for macOS systems
type MacDaemonInstaller struct {
	baseFolder string
}

func NewMacDaemonInstaller(baseFolder string) *MacDaemonInstaller {
	return &MacDaemonInstaller{baseFolder: baseFolder}
}

func (m *MacDaemonInstaller) CheckAndStopExistingService() error {
	color.Yellow.Println("🔍 Checking if service is running...")
	cmd := exec.Command("launchctl", "list", "xyz.shelltime.daemon")
	if err := cmd.Run(); err == nil {
		color.Yellow.Println("🛑 Stopping existing service...")
		if err := exec.Command("launchctl", "unload", "/Library/LaunchDaemons/xyz.shelltime.daemon.plist").Run(); err != nil {
			return fmt.Errorf("failed to stop existing service: %w", err)
		}
	}
	return nil
}

func (m *MacDaemonInstaller) InstallService() error {
	daemonPath := filepath.Join(m.baseFolder, "daemon")
	// Create daemon directory if not exists
	if err := os.MkdirAll(daemonPath, 0755); err != nil {
		return fmt.Errorf("failed to create daemon directory: %w", err)
	}

	plistPath := filepath.Join(daemonPath, "xyz.shelltime.daemon.plist")
	if _, err := os.Stat(plistPath); err == nil {
		if err := os.Remove(plistPath); err != nil {
			return fmt.Errorf("failed to remove existing plist file: %w", err)
		}
	}

	if err := os.WriteFile(plistPath, daemonMacServiceDesc, 0644); err != nil {
		return fmt.Errorf("failed to write plist file: %w", err)
	}
	return nil
}

func (m *MacDaemonInstaller) RegisterService() error {
	plistPath := "/Library/LaunchDaemons/xyz.shelltime.daemon.plist"
	if _, err := os.Stat(plistPath); err != nil {
		sourceFile := filepath.Join(m.baseFolder, "daemon/xyz.shelltime.daemon.plist")
		if err := os.Symlink(sourceFile, plistPath); err != nil {
			return fmt.Errorf("failed to create plist symlink: %w", err)
		}
	}
	return nil
}

func (m *MacDaemonInstaller) StartService() error {
	color.Yellow.Println("🚀 Starting service...")
	if err := exec.Command("launchctl", "load", "/Library/LaunchDaemons/xyz.shelltime.daemon.plist").Run(); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}
	return nil
}
