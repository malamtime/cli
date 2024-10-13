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
	},
	Action: commandGC,
}

func commandGC(c *cli.Context) error {
	storageFolder := os.ExpandEnv("$HOME/" + model.COMMAND_BASE_STORAGE_FOLDER)
	if _, err := os.Stat(storageFolder); os.IsNotExist(err) {
		return nil
	}

	if c.Bool("withLog") {
		logFile := os.ExpandEnv("$HOME/.malamtime/log.log")
		if err := os.Remove(logFile); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove log file: %v", err)
		}
	}

	// only can setup logger after the log file clean
	SetupLogger("$HOME/" + model.COMMAND_BASE_STORAGE_FOLDER)
	defer CloseLogger()

	commandsFolder := os.ExpandEnv("$HOME/" + model.COMMAND_STORAGE_FOLDER)
	if _, err := os.Stat(commandsFolder); os.IsNotExist(err) {
		return nil
	}

	lastCursor, err := getLastCursor()
	if err != nil {
		return err
	}

	postCommands, postCount, err := getPostCommands()
	if err != nil {
		return err
	}

	if postCount == 0 {
		logrus.Traceln("no post commands need to be clean")
		return nil
	}

	preTree, err := getPreCommands()
	if err != nil {
		return err
	}

	newPreCommandList := make([]*model.Command, 0)
	newPostCommandList := make([][]byte, 0)

	for _, row := range postCommands {
		parts := bytes.Split(row, []byte{model.SEPARATOR})
		if len(parts) < 2 {
			continue
		}
		cmd := new(model.Command)
		recordingTime, err := cmd.FromLine(string(parts[0]))
		if err != nil {
			logrus.Errorf("Failed to parse command: %v", err)
			continue
		}
		if recordingTime.Before(lastCursor) {
			continue
		}
		key := cmd.GetUniqueKey()
		preList, ok := preTree[key]
		if !ok {
			logrus.Warnln("could not find the match cmd in preFile", key)
			continue
		}
		closest := cmd.FindClosestCommand(preList)
		if closest == nil {
			logrus.Warnln("could not match the closest cmd in preList", key, len(preList))
			continue
		}
		newPreCommandList = append(newPreCommandList, closest)
		newPostCommandList = append(newPostCommandList, row)
	}

	sort.Slice(newPreCommandList, func(i, j int) bool {
		return newPreCommandList[i].RecordingTime.Before(newPreCommandList[j].RecordingTime)
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

	postFileContent := bytes.Join(newPostCommandList, []byte("\n"))
	if err := os.WriteFile(originalPostFile, postFileContent, 0644); err != nil {
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
