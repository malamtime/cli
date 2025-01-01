package model

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"go.opentelemetry.io/otel"
)

type GinGraphQLContextType struct {
	IP     string
	UserID int
}

var commitID string
var modelTracer = otel.Tracer("model")

const MAX_BUFFER_SIZE = 512 * 1024 // 512Kb

func InjectVar(commitId string) {
	commitID = commitId
}

// SudoGetBaseFolder will return the first matched `~/.shelltime/` folder
func SudoGetBaseFolder() (string, string, error) {
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
			entryInfo, err := entry.Info()
			if err != nil {
				return "", "", nil
			}

			isSymbolLink := entryInfo.Mode()&os.ModeSymlink != 0

			if !entry.IsDir() && !isSymbolLink {
				continue
			}

			shelltimePath := filepath.Join(basePath, entry.Name(), ".shelltime", "bin")
			if _, err := os.Stat(shelltimePath); err == nil {
				foundUser = entry.Name()
				break
			}
		}
		if foundUser != "" {
			break
		}
	}

	baseFolder, err := SudoGetUserBaseFolder(foundUser)
	return baseFolder, foundUser, err
}

func SudoGetUserBaseFolder(username string) (string, error) {
	homeAbsolutePrefix := ""
	if runtime.GOOS == "linux" {
		homeAbsolutePrefix = "/home"
	} else if runtime.GOOS == "darwin" {
		homeAbsolutePrefix = "/Users"
	}

	if username == "" && runtime.GOOS == "linux" {
		shelltimePath := filepath.Join("/root", ".shelltime", "bin")
		if _, err := os.Stat(shelltimePath); err == nil {
			username = "root"
		}
	}

	if username == "" {
		return "", fmt.Errorf("could not find any user with ~/.shelltime/bin directory")
	}

	if username == "root" && runtime.GOOS == "linux" {
		return filepath.Join("/root", ".shelltime"), nil
	}

	return filepath.Join(homeAbsolutePrefix, username, ".shelltime"), nil
}
