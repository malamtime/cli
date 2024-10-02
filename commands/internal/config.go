package internal

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

func ReadConfigFile() (config MalamTimeConfig, err error) {
	configFile := os.ExpandEnv("$HOME/.malamtime/config.toml")
	existingConfig, err := os.ReadFile(configFile)
	if err != nil {
		err = fmt.Errorf("failed to read config file: %w", err)
		return
	}

	err = toml.Unmarshal(existingConfig, &config)
	if err != nil {
		err = fmt.Errorf("failed to parse config file: %w", err)
		return
	}

	return
}
