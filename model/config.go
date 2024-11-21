package model

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

var UserShellTimeConfig ShellTimeConfig

//go:generate mockery --name ConfigService
type ConfigService interface {
	ReadConfigFile() (ShellTimeConfig, error)
}

type configService struct {
	configFilePath string
}

func NewConfigService(configFilePath string) ConfigService {
	return &configService{
		configFilePath: configFilePath,
	}
}

func (cs *configService) ReadConfigFile() (config ShellTimeConfig, err error) {
	configFile := cs.configFilePath
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

	// default 10 and at least 3 for performance reason
	if config.FlushCount == 0 {
		config.FlushCount = 10
	}
	if config.FlushCount < 3 {
		config.FlushCount = 3
	}

	if config.GCTime == 0 {
		config.GCTime = 14
	}

	truthy := true
	if config.DataMasking == nil {
		config.DataMasking = &truthy
	}
	UserShellTimeConfig = config
	return
}
