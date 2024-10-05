package commands

import "github.com/urfave/cli/v2"


var GCCommand *cli.Command = &cli.Command{
	Name:  "track",
	Usage: "track user commands",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name: "isUninstall",
			Usage: "if with uninstall parameter. will delete all logs, hooks, local data",
			Action: commandGC,
		}
	},
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

	dbFolder := filepath.Join(homeDir, ".malamtime", "db")
	if err := os.RemoveAll(dbFolder); err != nil {
		return fmt.Errorf("failed to remove db folder: %v", err)
	}

	if !c.Bool("isUninstall") {
		return nil
	}

	// TODO: delete $HOME/.config/malamtime/ folder

	return nil
}
