package model

import (
	"fmt"
	"runtime"
)

// DaemonInstaller interface defines methods for daemon installation across different platforms
type DaemonInstaller interface {
	CheckAndStopExistingService() error
	InstallService() error
	RegisterService() error
	StartService() error
}

// Factory function to create appropriate installer based on OS
func NewDaemonInstaller(baseFolder string) (DaemonInstaller, error) {
	switch runtime.GOOS {
	case "linux":
		return NewLinuxDaemonInstaller(baseFolder), nil
	case "darwin":
		return NewMacDaemonInstaller(baseFolder), nil
	default:
		return nil, fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

// func commandDaemonInstall(c *cli.Context) error {
//     // Check if running as root
//     if os.Geteuid() != 0 {
//         return fmt.Errorf("this command must be run as root (sudo shelltime daemon:install)")
//     }

//     installer, err := NewDaemonInstaller(sysDescFS)
//     if err != nil {
//         return err
//     }

//     if err := installer.CheckAndStopExistingService(); err != nil {
//         return err
//     }

//     realBin, err := getRealBinPath()
//     if err != nil {
//         return err
//     }

//     // ... (other setup code like checking for backup files)

//     if err := installer.InstallService(); err != nil {
//         return err
//     }

//     if err := installer.RegisterService(realBin); err != nil {
//         return err
//     }

//     if err := installer.StartService(); err != nil {
//         return err
//     }

//     color.Green.Println("✅ Daemon service has been installed and started successfully!")
//     color.Green.Println("💡 Your commands will now be automatically synced to shelltime.xyz")
//     return nil
// }
