package commands

import (
	"context"
	"fmt"
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
	SetupLogger(os.ExpandEnv("$HOME/" + model.COMMAND_BASE_STORAGE_FOLDER))
	defer CloseLogger()

	ctx := c.Context
	logrus.Trace(c.Args().First())
	config, err := configService.ReadConfigFile()
	if err != nil {
		logrus.Errorln(err)
		return err
	}

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
		Shell:     shell,
		SessionID: sessionId,
		Command:   cmdCommand,
		Hostname:  hostname,
		Username:  username,
		Time:      time.Now(),
		Phase:     model.CommandPhasePre,
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

	if cmdPhase == "post" {
		return trySyncLocalToServer(ctx, config)
	}
	return nil
}
func trySyncLocalToServer(ctx context.Context, config model.MalamTimeConfig) error {
	postFileContent, lineCount, err := model.GetPostCommands()
	if err != nil {
		return err
	}
	// if lineCount%config.FlushCount != 0 {
	// 	logrus.Traceln("Not enough records to sync, current count:", lineCount)
	// 	return nil
	// }

	if len(postFileContent) == 0 || lineCount == 0 {
		logrus.Traceln("Not enough records to sync, current count:", lineCount)
		return nil
	}

	cursor, err := model.GetLastCursor()
	if err != nil {
		return err
	}

	preFileTree, err := model.GetPreCommandsTree()
	if err != nil {
		return err
	}

	trackingData := make([]model.TrackingData, 0)

	var latestRecordingTime time.Time = cursor

	for _, line := range postFileContent {
		postCommand := new(model.Command)
		recordingTime, err := postCommand.FromLine(string(line))
		if err != nil {
			logrus.Errorln("Failed to parse post command:", err)
			continue
		}

		if recordingTime.Before(cursor) {
			continue
		}
		if recordingTime.After(latestRecordingTime) {
			latestRecordingTime = recordingTime
		}

		key := postCommand.GetUniqueKey()
		preCommands, ok := preFileTree[key]
		if !ok {
			continue
		}

		// here very sure the commandList are all elligable, so no need check here.
		closestPreCommand := postCommand.FindClosestCommand(preCommands, false)

		td := model.TrackingData{
			Shell:     postCommand.Shell,
			SessionID: postCommand.SessionID,
			Command:   postCommand.Command,
			Hostname:  postCommand.Hostname,
			Username:  postCommand.Username,
			EndTime:   postCommand.Time.Unix(),
			Result:    postCommand.Result,
		}

		if closestPreCommand != nil {
			td.StartTime = closestPreCommand.Time.Unix()
		}

		trackingData = append(trackingData, td)
	}

	if len(trackingData) == 0 {
		logrus.Traceln("no tracking data need to be sync")
		return nil
	}

	if len(trackingData) < config.FlushCount {
		logrus.Traceln("not enough data need to flush, abort. current is:", len(trackingData))
		return nil
	}

	err = model.SendLocalDataToServer(ctx, config, latestRecordingTime, trackingData)
	if err != nil {
		logrus.Errorln("Failed to send data to server:", err)
		return err
	}
	// TODO: update cursor
	return updateCursorToFile(latestRecordingTime)
}

func updateCursorToFile(latestRecordingTime time.Time) error {
	cursorFilePath := os.ExpandEnv(fmt.Sprintf("%s/%s", "$HOME", model.COMMAND_CURSOR_STORAGE_FILE))
	cursorFile, err := os.OpenFile(cursorFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		logrus.Errorln("Failed to open cursor file for writing:", err)
		return err
	}
	defer cursorFile.Close()

	_, err = cursorFile.WriteString(fmt.Sprintf("%d\n", latestRecordingTime.UnixNano()))
	if err != nil {
		logrus.Errorln("Failed to write to cursor file:", err)
		return err
	}
	return nil
}
