package commands

import (
	"context"
	"os"
	"time"

	"github.com/malamtime/cli/ent"
	"github.com/malamtime/cli/ent/command"
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
	if cmdPhase == "pre" {
		_, err = model.EntClient.Command.Create().
			SetShell(shell).
			SetSessionId(sessionId).
			SetCommand(cmdCommand).
			SetHostname(hostname).
			SetUsername(username).
			SetTime(time.Now()).
			SetPhase(command.PhasePre).
			SetSentToServer(false).
			Save(ctx)
	}

	if cmdPhase == "post" {
		cmd, err := model.EntClient.Command.Query().
			Where(
				command.Shell(shell),
				command.SessionIdEQ(sessionId),
				command.CommandEQ(cmdCommand),
				command.UsernameEQ(username),
				command.PhaseEQ(command.PhasePre),
				command.TimeGTE(time.Now().AddDate(0, 0, -10)),
			).
			Order(ent.Desc(command.FieldID)).
			First(ctx)
		if err != nil {
			logrus.Errorln("Failed to find matching command: ", err)
			return nil
		} else {
			_, err = cmd.Update().
				SetResult(result).
				SetEndTime(time.Now()).
				Save(ctx)
			if err != nil {
				logrus.Errorln("Failed to update command: ", err)
				return nil
			}
		}
	}

	if err != nil {
		logrus.Errorln(err)
		return err
	}

	return trySyncLocalToServer(ctx, config)
}

func trySyncLocalToServer(ctx context.Context, config model.MalamTimeConfig) error {
	count, err := model.EntClient.Command.Query().
		Where(command.SentToServerEQ(false), command.PhaseEQ(command.PhasePost)).
		Count(ctx)
	if err != nil {
		logrus.Errorln("Failed to get count of unsent commands:", err)
		return err
	}

	// do nothing if less than 10 records
	if count < config.FlushCount {
		return nil
	}

	commands, err := model.EntClient.Command.Query().
		Where(command.SentToServerEQ(false), command.PhaseEQ(command.PhasePost)).
		All(ctx)
	if err != nil {
		logrus.Errorln("Failed to retrieve unsent commands:", err)
		return err
	}

	trackingData := make([]model.TrackingData, len(commands))
	for i, cmd := range commands {
		trackingData[i] = model.TrackingData{
			Shell:     cmd.Shell,
			SessionID: cmd.SessionId,
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

	commandIds := make([]int, len(commands))
	for i, cmd := range commands {
		commandIds[i] = cmd.ID
	}

	_, err = model.EntClient.Command.Update().
		Where(
			command.SentToServerEQ(false),
			command.PhaseEQ(command.PhasePost),
			command.IDIn(commandIds...),
		).
		SetSentToServer(true).
		Save(ctx)
	if err != nil {
		logrus.Errorln("Failed to update sent status:", err)
		return err
	}

	gcStartTime := time.Now().AddDate(0, 0, -config.GCTime)
	_, err = model.EntClient.Command.Delete().
		Where(
			command.And(
				command.SentToServerEQ(true),
				command.PhaseEQ(command.PhasePost),
				command.TimeLT(gcStartTime),
			),
		).
		Exec(ctx)
	if err != nil {
		logrus.Errorln("Failed to delete old records:", err)
		return err
	}

	return nil
}
