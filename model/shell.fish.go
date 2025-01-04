// service/fish_hook_service.go
package model

import (
	"fmt"
	"os"
	"path/filepath"
)

type FishHookService struct {
	BaseHookService

	hookLines []string
}

func NewFishHookService() *FishHookService {
	sourceContent := os.ExpandEnv("$HOME/.shelltime/hooks/fish.fish")
	hookLines := []string{
		"fish_add_path $HOME/.shelltime/bin",
		fmt.Sprintf("source %s", sourceContent),
	}
	return &FishHookService{
		hookLines: hookLines,
	}
}

func (s *FishHookService) Uninstall() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(home, ".config", "fish", "config.fish")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil
	}

	// Backup the file
	if err := s.backupFile(configPath); err != nil {
		return err
	}
	return s.removeHookLines(configPath, s.hookLines)
}
