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
	UnregisterService() error
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
