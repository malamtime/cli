package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gookit/color"
	"github.com/malamtime/cli/model"
	"github.com/urfave/cli/v2"
)

var DaemonInstallCommand *cli.Command = &cli.Command{
	Name:   "daemon:install",
	Usage:  "Install shelltime daemon service",
	Action: commandDaemonInstall,
}

func commandDaemonInstall(c *cli.Context) error {
	color.Yellow.Println("⚠️ Warning: This daemon service is currently not ready for use. Please proceed with caution.")

	// Check if running as root
	if os.Geteuid() != 0 {
		return fmt.Errorf("this command must be run as root (sudo shelltime daemon:install)")
	}
	color.Yellow.Println("🔍 Detecting system architecture...")

	baseFolder, err := getBaseFolder()
	if err != nil {
		return err
	}

	installer, err := model.NewDaemonInstaller(baseFolder)
	if err != nil {
		return err
	}

	if err := installer.CheckAndStopExistingService(); err != nil {
		return err
	}

	// check latest file exist or not
	if _, err := os.Stat(filepath.Join(baseFolder, "bin/shelltime-daemon.bak")); err == nil {
		color.Yellow.Println("🔄 Found latest daemon file, restoring...")
		// try to remove old file
		_ = os.Remove(filepath.Join(baseFolder, "bin/shelltime-daemon"))
		// rename .bak to original
		if err := os.Rename(
			filepath.Join(baseFolder, "bin/shelltime-daemon.bak"),
			filepath.Join(baseFolder, "bin/shelltime-daemon"),
		); err != nil {
			return fmt.Errorf("failed to restore latest daemon: %w", err)
		}
	}

	// check shelltime-daemon
	if _, err := os.Stat(filepath.Join(baseFolder, "bin/shelltime-daemon")); err != nil {
		color.Yellow.Println("⚠️ shelltime-daemon not found, please reinstall the CLI first:")
		color.Yellow.Println("curl -sSL https://raw.githubusercontent.com/malamtime/installation/master/install.bash | bash")
		return nil
	}

	// Copy to final location
	binaryPath := "/usr/local/bin/shelltime-daemon"

	if _, err := os.Stat(binaryPath); err != nil {
		color.Yellow.Println("🔍 Creating daemon symlink...")
		if err := os.Symlink(filepath.Join(baseFolder, "bin/shelltime-daemon"), binaryPath); err != nil {
			return fmt.Errorf("failed to create daemon symlink: %w", err)
		}
	}

	if err := installer.InstallService(); err != nil {
		return err
	}

	if err := installer.RegisterService(); err != nil {
		return err
	}

	return installer.StartService()
}

// getBaseFolder will return the first matched `~/.shelltime/` folder
func getBaseFolder() (string, error) {
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
