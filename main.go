package main

import (
	"fmt"
	"os"

	"github.com/malamtime/cli/commands"
	"github.com/malamtime/cli/model"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "print the version",
	}

	cli.VersionPrinter = func(cCtx *cli.Context) {
		fmt.Fprintf(cCtx.App.Writer, "version=%s\n", cCtx.App.Version)
	}

	configFile := os.ExpandEnv(fmt.Sprintf("%s/%s/%s", "$HOME", model.COMMAND_BASE_STORAGE_FOLDER, "config.toml"))
	configService := model.NewConfigService(configFile)

	model.InjectVar(version)
	commands.InjectVar(version, configService)
	app := cli.NewApp()
	app.Name = "shelltime CLI"
	app.Description = "shelltime.xyz CLI for track DevOps works"
	app.Usage = "shelltime.xyz CLI for track DevOps works"
	app.Version = version
	app.Copyright = "Copyright (c) 2024 shelltime.xyz Team"
	app.Authors = []*cli.Author{
		{
			Name:  "shelltime.xyz Team",
			Email: "annatar.he+shelltime.xyz@gmail.com",
		},
	}
	app.Suggest = true
	app.HideVersion = false
	app.Metadata = map[string]interface{}{
		"version": version,
	}

	app.Commands = []*cli.Command{
		commands.AuthCommand,
		commands.TrackCommand,
		commands.GCCommand,
	}
	err := app.Run(os.Args)
	if err != nil {
		logrus.Errorln(err)
	}
	commands.CloseLogger()
}
