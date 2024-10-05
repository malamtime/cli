package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

var GCCommand *cli.Command = &cli.Command{
	Name:  "gc",
	Usage: "clean internal storage",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "database",
			Aliases: []string{"db"},
			Usage:   "will clean local temporary data. will loose some data",
		},
	},
	Action: commandGC,
}

func commandGC(c *cli.Context) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %v", err)
	}

	logFile := filepath.Join(homeDir, ".malamtime", "log.log")
	if err := os.Remove(logFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove log file: %v", err)
	}

	if c.Bool("database") {
		dbFolder := filepath.Join(homeDir, ".malamtime", "db")
		if err := os.RemoveAll(dbFolder); err != nil {
			return fmt.Errorf("failed to remove db folder: %v", err)
		}
	}

	// TODO: delete $HOME/.config/malamtime/ folder

	return nil
}
