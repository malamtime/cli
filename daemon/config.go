package daemon

import (
	"os"
	"path/filepath"

	"github.com/malamtime/cli/model"
	"gopkg.in/yaml.v3"
)

const (
	DefaultSocketPath = "/tmp/shelltime.sock"
	DefaultConfigPath = "/etc/shelltime/config.yml"
)

type DaemonConfig struct {
	SocketPath string `yaml:"socketPath"`
	SystemUser string `yaml:"sysUser"`
}

// ConfigService defines the interface for daemon configuration operations
type ConfigService interface {
	GetConfig() (*DaemonConfig, error)
	CreateDefault() (*DaemonConfig, error)
	GetUserConfig() (model.ConfigService, error)
}

// configService implements ConfigService
type configService struct {
	configPath string
}

// NewConfigService creates a new instance of ConfigService
func NewConfigService(configPath string) ConfigService {
	return &configService{
		configPath: configPath,
	}
}

// GetConfig reads and returns the daemon configuration
func (s *configService) GetConfig() (*DaemonConfig, error) {
	if _, err := os.Stat(s.configPath); os.IsNotExist(err) {
		return s.CreateDefault()
	}

	data, err := os.ReadFile(s.configPath)
	if err != nil {
		return nil, err
	}

	var config DaemonConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// CreateDefault creates and returns a default daemon configuration
func (s *configService) CreateDefault() (*DaemonConfig, error) {

	_, username, err := model.SudoGetBaseFolder()
	if err != nil {
		return nil, err
	}

	config := &DaemonConfig{
		SocketPath: DefaultSocketPath,
		SystemUser: username,
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(filepath.Dir(s.configPath), 0755)
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(s.configPath, data, 0644)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (s *configService) GetUserConfig() (model.ConfigService, error) {
	c, err := s.GetConfig()
	if err != nil {
		return nil, err
	}

	baseFolder, err := model.SudoGetUserBaseFolder(c.SystemUser)
	if err != nil {
		return nil, err
	}

	return model.NewConfigService(
		filepath.Join(baseFolder, "config.toml"),
	), nil
}
