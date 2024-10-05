package main

import (
	"os"

	"github.com/malamtime/cli/commands"
	"github.com/malamtime/cli/model"
	"github.com/urfave/cli/v2"
)

var GitCommit string

func main() {
	model.InjectVar(GitCommit)
	commands.InjectVar(GitCommit)
	commands.SetupLogger()
	defer commands.CloseLogger()
	app := &cli.App{
		Name:        "MalamTime CLI",
		Description: "MalamTime CLI for track DevOps works",
		Usage:       "MalamTime CLI for track DevOps works",
		Version:     GitCommit,
		Commands: []*cli.Command{
			commands.AuthCommand,
			commands.TrackCommand,
		},
	}

	app.Run(os.Args)
	// if err != nil {
	// 	color.Red.Println(err)
	// }
}
