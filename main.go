package main

import (
	"fmt"
	"os"

	"github.com/malamtime/cli/commands"
	"github.com/malamtime/cli/model"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var GitCommit string

func main() {
	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "print the version",
	}

	cli.VersionPrinter = func(cCtx *cli.Context) {
		fmt.Fprintf(cCtx.App.Writer, "version=%s\n", cCtx.App.Version)
	}

	configFile := os.ExpandEnv("$HOME/.malamtime/config.toml")
	configService := model.NewConfigService(configFile)

	model.InjectVar(GitCommit)
	commands.InjectVar(GitCommit, configService)
	app := cli.NewApp()
	app.Name = "shelltime CLI"
	app.Description = "shelltime.xyz CLI for track DevOps works"
	app.Usage = "shelltime.xyz CLI for track DevOps works"
	app.Version = GitCommit
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
		"version": GitCommit,
	}

	app.Commands = []*cli.Command{
		commands.AuthCommand,
		commands.TrackCommand,
		commands.GCCommand,
		{
			Name:    "version",
			Aliases: []string{"v"},
			Action: func(ctx *cli.Context) error {
				fmt.Println(GitCommit)
				return nil
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		logrus.Errorln(err)
	}
	commands.CloseLogger()
}
