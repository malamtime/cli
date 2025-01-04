// service/zsh_hook_service.go
package model

import (
	"fmt"
	"os"
	"path/filepath"
)

type ZshHookService struct {
	BaseHookService
	hookLines []string
}

func NewZshHookService() *ZshHookService {
	sourceContent := os.ExpandEnv("$HOME/.shelltime/hooks/zsh.zsh")
	return &ZshHookService{
		hookLines: []string{
			`export PATH="$HOME/.shelltime/bin:$PATH"`,
			fmt.Sprintf("source %s", sourceContent),
		},
	}
}

func (s *ZshHookService) Uninstall() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	rcPath := filepath.Join(home, ".zshrc")
	if _, err := os.Stat(rcPath); os.IsNotExist(err) {
		return nil
	}

	// Backup the file
	if err := s.backupFile(rcPath); err != nil {
		return err
	}

	return s.removeHookLines(rcPath, s.hookLines)
}
