package commands

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
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
	logrus.Trace(c.Args().First())
	config, err := configService.ReadConfigFile()
	if err != nil {
		logrus.Errorln(err)
		return err
	}

	// model.InitDB()
	// defer model.Clean()
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

	return trySyncLocalToServer(ctx, config)
}

func getLastCursor() (cursorTime time.Time, err error) {
	cursorFilePath := os.ExpandEnv(fmt.Sprintf("%s/%s", "$HOME", model.COMMAND_CURSOR_STORAGE_FILE))
	cursorFile, err := os.Open(cursorFilePath)
	if err != nil {
		logrus.Errorln("Failed to open cursor file:", err)
		return
	}
	defer cursorFile.Close()

	scanner := bufio.NewScanner(cursorFile)
	var lastLine string
	for scanner.Scan() {
		lastLine = scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		logrus.Errorln("Error reading cursor file:", err)
		return cursorTime, err
	}

	cursor, err := strconv.Atoi(lastLine)
	if err != nil {
		logrus.Errorln("Failed to parse cursor value:", err)
		return
	}
	cursorTime = time.Unix(0, int64(cursor))
	return
}

func getPostCommands() ([][]byte, int, error) {
	postFilePath := os.ExpandEnv(fmt.Sprintf("%s/%s", "$HOME", model.COMMAND_POST_STORAGE_FILE))
	postFileHandler, err := os.Open(postFilePath)
	if err != nil {
		logrus.Errorln("Failed to open file:", err)
		return nil, 0, err
	}
	defer postFileHandler.Close()

	var fileCount [][]byte

	lineCount := 0
	postFileScanner := bufio.NewScanner(postFileHandler)
	for postFileScanner.Scan() {
		fileCount = append(fileCount, postFileScanner.Bytes())
		lineCount++
	}
	if err := postFileScanner.Err(); err != nil {
		logrus.Errorln("Error reading file:", err)
		return nil, 0, err
	}

	return fileCount, lineCount, nil
}

// key: ${shell}|${sessionID}|${command}|${username}
// value: model.Command
type preCommandTree map[string][]*model.Command

func getPreCommands() (result preCommandTree, err error) {
	preFilePath := os.ExpandEnv(fmt.Sprintf("%s/%s", "$HOME", model.COMMAND_PRE_STORAGE_FILE))
	preFileHandler, err := os.Open(preFilePath)
	if err != nil {
		logrus.Errorln("Failed to open pre-command file:", err)
		return nil, err
	}
	defer preFileHandler.Close()

	preFileScanner := bufio.NewScanner(preFileHandler)
	for preFileScanner.Scan() {
		line := preFileScanner.Text()
		var cmd *model.Command

		_, err := cmd.FromLine(line)
		if err != nil {
			logrus.Errorln("Invalid line parse in pre-command file:", line, err)
			continue
		}

		key := cmd.GetUniqueKey()
		result[key] = append(result[key], cmd)
	}

	if err := preFileScanner.Err(); err != nil {
		logrus.Errorln("Error reading pre-command file:", err)
		return nil, err
	}

	return result, nil
}

func trySyncLocalToServer(ctx context.Context, config model.MalamTimeConfig) error {
	postFileContent, lineCount, err := getPostCommands()
	if err != nil {
		return err
	}
	if lineCount%config.FlushCount != 0 {
		logrus.Traceln("Not enough records to sync, current count:", lineCount)
		return nil
	}

	if len(postFileContent) == 0 || lineCount == 0 {
		logrus.Traceln("Not enough records to sync, current count:", lineCount)
		return nil
	}

	cursor, err := getLastCursor()
	if err != nil {
		return err
	}

	preFileTree, err := getPreCommands()
	if err != nil {
		return err
	}

	trackingData := make([]model.TrackingData, 0)

	var latestRecordingTime time.Time = cursor

	for _, line := range postFileContent {
		var postCommand *model.Command
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

		var closestPreCommand *model.Command
		minTimeDiff := int64(^uint64(0) >> 1) // Max int64 value

		for _, preCommand := range preCommands {
			timeDiff := postCommand.Time.Unix() - preCommand.Time.Unix()
			if timeDiff >= 0 && timeDiff < minTimeDiff {
				minTimeDiff = timeDiff
				closestPreCommand = preCommand
			}
		}

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

	err = model.SendLocalDataToServer(ctx, config, trackingData)
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
