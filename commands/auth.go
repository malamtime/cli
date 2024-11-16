package commands

import (
	"fmt"
	"os"

	"github.com/gookit/color"
	"github.com/malamtime/cli/model"
	"github.com/pelletier/go-toml/v2"
	"github.com/urfave/cli/v2"
)

var AuthCommand *cli.Command = &cli.Command{
	Name:  "init",
	Usage: "init your shelltime.xyz config",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "token",
			Aliases:  []string{"t"},
			Usage:    "Authentication token",
			Required: true,
		},
	},
	Action: commandAuth,
}

func commandAuth(c *cli.Context) error {
	configDir := os.ExpandEnv("$HOME/" + model.COMMAND_BASE_STORAGE_FOLDER)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err = os.Mkdir(configDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}
	}
	SetupLogger(configDir)

	var config model.ShellTimeConfig
	configFile := configDir + "/config.toml"
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		content, err := toml.Marshal(model.DefaultConfig)
		if err != nil {
			return fmt.Errorf("failed to marshal default config: %w", err)
		}
		err = os.WriteFile(configFile, content, 0644)
		if err != nil {
			return fmt.Errorf("failed to create config file: %w", err)
		}
		config = model.DefaultConfig
	} else {
		existingConfig, err := configService.ReadConfigFile()
		if err != nil {
			return err
		}
		config = existingConfig
	}

	newToken := c.String("token")
	config.Token = newToken
	content, err := toml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	err = os.WriteFile(configFile, content, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	color.Green.Println(" ✅ config file created")
	return nil
}
