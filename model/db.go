package model

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	SEPARATOR = byte('\t')
)

var (
	COMMAND_BASE_STORAGE_FOLDER = ".malamtime"
	COMMAND_STORAGE_FOLDER      = COMMAND_BASE_STORAGE_FOLDER + "/commands"
	COMMAND_PRE_STORAGE_FILE    = COMMAND_STORAGE_FOLDER + "/pre.txt"
	COMMAND_POST_STORAGE_FILE   = COMMAND_STORAGE_FOLDER + "/post.txt"
	COMMAND_CURSOR_STORAGE_FILE = COMMAND_STORAGE_FOLDER + "/cursor.txt"
)

func InitFolder(baseFolder string) {
	if baseFolder != "" {
		COMMAND_BASE_STORAGE_FOLDER = fmt.Sprintf(".malamtime-%s", baseFolder)
	}

	COMMAND_STORAGE_FOLDER = COMMAND_BASE_STORAGE_FOLDER + "/commands"
	COMMAND_PRE_STORAGE_FILE = COMMAND_STORAGE_FOLDER + "/pre.txt"
	COMMAND_POST_STORAGE_FILE = COMMAND_STORAGE_FOLDER + "/post.txt"
	COMMAND_CURSOR_STORAGE_FILE = COMMAND_STORAGE_FOLDER + "/cursor.txt"
}

// key: ${shell}|${sessionID}|${command}|${username}
// value: model.Command
type preCommandTree map[string][]*Command

func GetPreCommandsTree() (result preCommandTree, err error) {
	preFilePath := os.ExpandEnv(fmt.Sprintf("%s/%s", "$HOME", COMMAND_PRE_STORAGE_FILE))
	preFileHandler, err := os.Open(preFilePath)
	if err != nil {
		logrus.Errorln("Failed to open pre-command file:", err)
		return nil, err
	}
	defer preFileHandler.Close()

	result = make(preCommandTree)
	preFileScanner := bufio.NewScanner(preFileHandler)
	for preFileScanner.Scan() {
		line := preFileScanner.Text()
		cmd := new(Command)

		_, err := cmd.FromLine(line)
		if err != nil {
			logrus.Errorln("Invalid line parse in pre-command file:", line, err)
			continue
		}

		key := cmd.GetUniqueKey()
		if len(result[key]) == 0 {
			result[key] = []*Command{cmd}
		} else {
			result[key] = append(result[key], cmd)
		}
	}

	if err := preFileScanner.Err(); err != nil {
		logrus.Errorln("Error reading pre-command file:", err)
		return nil, err
	}

	return result, nil
}

func GetPreCommands() ([]*Command, error) {
	preFilePath := os.ExpandEnv(fmt.Sprintf("%s/%s", "$HOME", COMMAND_PRE_STORAGE_FILE))
	preFileHandler, err := os.Open(preFilePath)
	if err != nil {
		logrus.Errorln("Failed to open pre-command file:", err)
		return nil, err
	}
	defer preFileHandler.Close()

	result := make([]*Command, 0)
	preFileScanner := bufio.NewScanner(preFileHandler)
	for preFileScanner.Scan() {
		line := preFileScanner.Text()
		cmd := new(Command)

		_, err := cmd.FromLine(line)
		if err != nil {
			logrus.Errorln("Invalid line parse in pre-command file:", line, err)
			continue
		}

		result = append(result, cmd)
	}

	if err := preFileScanner.Err(); err != nil {
		logrus.Errorln("Error reading pre-command file:", err)
		return nil, err
	}

	return result, nil
}

func GetLastCursor() (cursorTime time.Time, err error) {
	cursorFilePath := os.ExpandEnv(fmt.Sprintf("%s/%s", "$HOME", COMMAND_CURSOR_STORAGE_FILE))
	cursorFile, err := os.Open(cursorFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			cursorTime = time.Time{}
			err = nil
			return
		}
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
	// if not data exists, just use time.Zero
	if lastLine == "" {
		return
	}

	cursor, err := strconv.Atoi(lastLine)
	if err != nil {
		logrus.Errorln("Failed to parse cursor value:", err)
		return
	}
	cursorTime = time.Unix(0, int64(cursor))
	return
}

func GetPostCommands() ([][]byte, int, error) {
	postFilePath := os.ExpandEnv(fmt.Sprintf("%s/%s", "$HOME", COMMAND_POST_STORAGE_FILE))
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
