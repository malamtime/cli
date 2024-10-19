package main

import (
	"os"

	"github.com/malamtime/cli/commands"
	"github.com/malamtime/cli/model"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var GitCommit string

func main() {
	configFile := os.ExpandEnv("$HOME/.malamtime/config.toml")
	configService := model.NewConfigService(configFile)

	model.InjectVar(GitCommit)
	commands.InjectVar(GitCommit, configService)
	app := &cli.App{
		Name:        "MalamTime CLI",
		Description: "MalamTime CLI for track DevOps works",
		Usage:       "MalamTime CLI for track DevOps works",
		Version:     GitCommit,
		Commands: []*cli.Command{
			commands.AuthCommand,
			commands.TrackCommand,
			commands.GCCommand,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.Errorln(err)
	}
	// every commands init logger, and here to close
	commands.CloseLogger()
}
