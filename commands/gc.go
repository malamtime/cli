package commands

import (
	"bytes"
	"fmt"
	"os"
	"sort"

	"github.com/malamtime/cli/model"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var GCCommand *cli.Command = &cli.Command{
	Name:  "gc",
	Usage: "clean internal storage",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "withLog",
			Aliases: []string{"wl"},
			Usage:   "clean the log file",
		},
		&cli.BoolFlag{
			Name:        "skipLogCreation",
			Aliases:     []string{"slc"},
			DefaultText: "false",
			Usage:       "skip log file creation",
		},
	},
	Action: commandGC,
}

func commandGC(c *cli.Context) error {
	storageFolder := os.ExpandEnv("$HOME/" + model.COMMAND_BASE_STORAGE_FOLDER)
	if _, err := os.Stat(storageFolder); os.IsNotExist(err) {
		return nil
	}

	if c.Bool("withLog") {
		logFile := os.ExpandEnv("$HOME/" + model.COMMAND_BASE_STORAGE_FOLDER + "/log.log")
		if err := os.Remove(logFile); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove log file: %v", err)
		}
	}

	if !c.Bool("skipLogCreation") {
		// only can setup logger after the log file clean
		SetupLogger(storageFolder)
		defer CloseLogger()
	}

	commandsFolder := os.ExpandEnv("$HOME/" + model.COMMAND_STORAGE_FOLDER)
	if _, err := os.Stat(commandsFolder); os.IsNotExist(err) {
		return nil
	}

	lastCursor, err := model.GetLastCursor()
	if err != nil {
		return err
	}

	postCommandsRaw, postCount, err := model.GetPostCommands()
	if err != nil {
		return err
	}

	postCommands := make([]*model.Command, len(postCommandsRaw))
	for i, row := range postCommandsRaw {
		cmd := new(model.Command)
		_, err := cmd.FromLine(string(row))
		if err != nil {
			err = fmt.Errorf("failed to parse command from line: %v", err)
			logrus.Warnln(err)
			continue
		}
		postCommands[i] = cmd
	}

	if postCount == 0 {
		logrus.Traceln("no post commands need to be clean")
		return nil
	}

	preCommands, err := model.GetPreCommands()
	if err != nil {
		return err
	}

	newPreCommandList := make([]*model.Command, 0)
	newPostCommandList := make([]*model.Command, 0)

	// save all the data that before cursor
	for _, cmd := range postCommands {
		if cmd.RecordingTime.After(lastCursor) {
			newPostCommandList = append(newPostCommandList, cmd)
		}
	}

	// If there is no end, it should be kept. For example, if one tab opened a webpack dev server and the user opened another tab, we should keep the previous pre
	for _, row := range preCommands {
		recordingTime := row.RecordingTime
		// If it's data after the cursor, save it without thinking
		if recordingTime.After(lastCursor) {
			newPreCommandList = append(newPreCommandList, row)
			continue
		}

		// if the closest node not found, prohaps the pre command not finished yet. save the pre command anyway
		closestNode := row.FindClosestCommand(postCommands, true)
		if closestNode == nil || closestNode.IsNil() {
			newPreCommandList = append(newPreCommandList, row)
		}
	}

	sort.Slice(newPreCommandList, func(i, j int) bool {
		return newPreCommandList[i].
			RecordingTime.
			Before(
				newPreCommandList[j].RecordingTime,
			)
	})

	sort.Slice(newPostCommandList, func(i, j int) bool {
		return newPostCommandList[i].
			RecordingTime.
			Before(
				newPostCommandList[j].RecordingTime,
			)
	})

	originalPreFile := os.ExpandEnv("$HOME/" + model.COMMAND_PRE_STORAGE_FILE)
	originalPostFile := os.ExpandEnv("$HOME/" + model.COMMAND_POST_STORAGE_FILE)
	originalCursorFile := os.ExpandEnv("$HOME/" + model.COMMAND_CURSOR_STORAGE_FILE)

	preBackupFile := originalPreFile + ".bak"
	postBackupFile := originalPostFile + ".bak"
	cursorBackupFile := originalCursorFile + ".bak"

	if err := os.Rename(originalPreFile, preBackupFile); err != nil {
		err = fmt.Errorf("failed to backup PRE_FILE: %v", err)
		logrus.Warnln(err)
		return err
	}
	if err := os.Rename(originalPostFile, postBackupFile); err != nil {
		err = fmt.Errorf("failed to backup POST_FILE: %v", err)
		logrus.Warnln(err)
		return err
	}
	if err := os.Rename(originalCursorFile, cursorBackupFile); err != nil {
		err = fmt.Errorf("failed to backup CURSOR_FILE: %v", err)
		logrus.Warnln(err)
		return err
	}

	preFileContent := bytes.Buffer{}
	for _, cmd := range newPreCommandList {
		line, err := cmd.ToLine(cmd.RecordingTime)
		if err != nil {
			return fmt.Errorf("failed to convert command to line: %v", err)
		}
		preFileContent.Write(line)
		preFileContent.Write([]byte{'\n'})
	}
	if err := os.WriteFile(originalPreFile, preFileContent.Bytes(), 0644); err != nil {
		err = fmt.Errorf("failed to write new PRE_FILE: %v", err)
		logrus.Warnln(err)
		return err
	}

	postFileContent := bytes.Buffer{}
	for _, cmd := range newPostCommandList {
		line, err := cmd.ToLine(cmd.RecordingTime)
		if err != nil {
			return fmt.Errorf("failed to convert command to line: %v", err)
		}
		postFileContent.Write(line)
		postFileContent.Write([]byte{'\n'})
	}

	if err := os.WriteFile(originalPostFile, postFileContent.Bytes(), 0644); err != nil {
		err = fmt.Errorf("failed to write new POST_FILE: %v", err)
		logrus.Warnln(err)
		return err
	}

	lastCursorNano := lastCursor.UnixNano()
	lastCursorBytes := []byte(fmt.Sprintf("%d", lastCursorNano))
	if err := os.WriteFile(originalCursorFile, lastCursorBytes, 0644); err != nil {
		err = fmt.Errorf("failed to write new CURSOR_FILE: %v", err)
		logrus.Warnln(err)
		return err
	}

	// TODO: delete $HOME/.config/malamtime/ folder

	return nil
}
