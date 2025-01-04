package commands

import (
	"fmt"
	"os"

	"github.com/gookit/color"
	"github.com/malamtime/cli/model"
	"github.com/urfave/cli/v2"
)

var DaemonUninstallCommand = &cli.Command{
	Name:   "uninstall",
	Usage:  "Uninstall the shelltime daemon service",
	Action: commandDaemonUninstall,
}

func commandDaemonUninstall(c *cli.Context) error {
	// Check if running as root
	if os.Geteuid() != 0 {
		return fmt.Errorf("this command must be run as root (sudo shelltime daemon uninstall)")
	}

	color.Yellow.Println("🔍 Starting daemon service uninstallation...")

	// TODO: the username is not stable in multiple user system
	baseFolder, username, err := model.SudoGetBaseFolder()
	if err != nil {
		return err
	}

	installer, err := model.NewDaemonInstaller(baseFolder, username)
	if err != nil {
		return err
	}

	// Unregister and remove the service
	if err := installer.UnregisterService(); err != nil {
		return fmt.Errorf("failed to unregister service: %w", err)
	}

	// Remove symlink from /usr/local/bin
	binaryPath := "/usr/local/bin/shelltime-daemon"
	if _, err := os.Stat(binaryPath); err == nil {
		color.Yellow.Println("🗑 Removing daemon symlink...")
		if err := os.Remove(binaryPath); err != nil {
			return fmt.Errorf("failed to remove daemon symlink: %w", err)
		}
	}

	color.Green.Println("✅ Daemon service has been successfully uninstalled!")
	// color.Yellow.Println("ℹ️  Note: Your commands will now be synced to shelltime.xyz on the next login")
	return nil
}
