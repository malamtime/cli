package commands

import (
	"context"
	"os"
	"time"

	"github.com/malamtime/cli/model"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var TrackCommand *cli.Command = &cli.Command{
	Name:  "track",
	Usage: "track user commands",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "shell",
			Aliases: []string{"s"},
			Value:   "",
			Usage:   "the shell that user use",
		},
		&cli.Int64Flag{
			Name:    "sessionId",
			Aliases: []string{"id"},
			Value:   0,
			Usage:   "unix timestamp of the session",
		},
		&cli.StringFlag{
			Name:    "command",
			Aliases: []string{"cmd"},
			Value:   "",
			Usage:   "command that user executed",
		},
		&cli.StringFlag{
			Name:    "phase",
			Aliases: []string{"p"},
			Usage:   "Phase: pre, post",
		},
		&cli.IntFlag{
			Name:    "result",
			Aliases: []string{"r"},
			Usage:   "Exit code of last command",
		},
	},
	Action: commandTrack,
	OnUsageError: func(cCtx *cli.Context, err error, isSubcommand bool) error {
		return nil
	},
}

func commandTrack(c *cli.Context) error {
	ctx := c.Context
	logrus.Info(c.Args().First())
	config, err := model.ReadConfigFile()
	if err != nil {
		logrus.Errorln(err)
		return err
	}

	model.InitDB()
	defer model.Clean()
	hostname, err := os.Hostname()
	if err != nil {
		logrus.Errorln(err)
		return err
	}

	username := os.Getenv("USER")

	shell := c.String("shell")
	sessionId := c.Int64("sessionId")
	cmdCommand := c.String("command")
	cmdPhase := c.String("phase")
	result := c.Int("result")

	instance := &model.Command{
		Shell:        shell,
		SessionID:    sessionId,
		Command:      cmdCommand,
		Hostname:     hostname,
		Username:     username,
		Time:         time.Now(),
		Phase:        model.CommandPhasePre,
		SentToServer: false,
	}

	if cmdPhase == "pre" {
		err = instance.DoSavePre()
	}
	if cmdPhase == "post" {
		err = instance.DoUpdate(result)
	}
	if err != nil {
		logrus.Errorln(err)
		return err
	}

	return trySyncLocalToServer(ctx, config)
}

func trySyncLocalToServer(ctx context.Context, config model.MalamTimeConfig) error {
	count, err := model.GetArchievedCount()
	if err != nil {
		logrus.Errorln("Failed to get count of unsent commands:", err)
		return err
	}

	// do nothing if less than 10 records
	if count < config.FlushCount {
		return nil
	}

	keys, commands, err := model.GetArchivedList(count)
	if err != nil {
		logrus.Errorln("Failed to retrieve unsent commands:", err)
		return err
	}

	trackingData := make([]model.TrackingData, len(commands))
	for i, cmd := range commands {
		trackingData[i] = model.TrackingData{
			Shell:     cmd.Shell,
			SessionID: cmd.SessionID,
			Command:   cmd.Command,
			Hostname:  cmd.Hostname,
			Username:  cmd.Username,
			StartTime: cmd.Time.Unix(),
			EndTime:   cmd.EndTime.Unix(),
			Result:    cmd.Result,
		}
	}

	err = model.SendLocalDataToServer(ctx, config, trackingData)
	if err != nil {
		logrus.Errorln("Failed to send data to server:", err)
		return err
	}

	err = model.CleanArchievedData(keys)
	if err != nil {
		logrus.Errorln("Failed to delete local archived data:", err)
		return err
	}
	return nil
}
