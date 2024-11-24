package model

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	SEPARATOR = byte('\t')
)

var (
	COMMAND_BASE_STORAGE_FOLDER = ".shelltime"
	COMMAND_STORAGE_FOLDER      = COMMAND_BASE_STORAGE_FOLDER + "/commands"
	COMMAND_PRE_STORAGE_FILE    = COMMAND_STORAGE_FOLDER + "/pre.txt"
	COMMAND_POST_STORAGE_FILE   = COMMAND_STORAGE_FOLDER + "/post.txt"
	COMMAND_CURSOR_STORAGE_FILE = COMMAND_STORAGE_FOLDER + "/cursor.txt"
)

func InitFolder(baseFolder string) {
	if baseFolder != "" {
		COMMAND_BASE_STORAGE_FOLDER = fmt.Sprintf(".shelltime-%s", baseFolder)
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

	fileContentRaw, err := io.ReadAll(preFileHandler)
	if err != nil {
		logrus.Errorln("Error reading pre-command file:", err)
		return nil, err
	}

	result = make(preCommandTree)

	fileContent := bytes.Split(fileContentRaw, []byte("\n"))

	for _, row := range fileContent {
		if len(row) == 0 {
			continue
		}
		line := string(row)
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

	fileContentRow, err := io.ReadAll(preFileHandler)
	if err != nil {
		logrus.Errorln("Error reading file:", err)
		return nil, err
	}

	fileContent := bytes.Split(fileContentRow, []byte("\n"))
	// Remove empty lines from fileContent
	nonEmptyContent := make([][]byte, 0)
	for _, line := range fileContent {
		if len(line) > 0 {
			nonEmptyContent = append(nonEmptyContent, line)
		}
	}
	fileContent = nonEmptyContent

	result := make([]*Command, 0)
	for _, row := range fileContent {
		line := string(row)
		cmd := new(Command)

		_, err := cmd.FromLine(line)
		if err != nil {
			logrus.Errorln("Invalid line parse in pre-command file:", line, err)
			continue
		}

		result = append(result, cmd)
	}

	return result, nil
}

func GetLastCursor() (cursorTime time.Time, noCursorExist bool, err error) {
	noCursorExist = false
	cursorFilePath := os.ExpandEnv(fmt.Sprintf("%s/%s", "$HOME", COMMAND_CURSOR_STORAGE_FILE))
	cursorFile, err := os.Open(cursorFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			cursorTime = time.Time{}
			noCursorExist = true
			err = nil
			return
		}
		logrus.Errorln("Failed to open cursor file:", err)
		return
	}
	defer cursorFile.Close()

	fileContent, err := io.ReadAll(cursorFile)

	if err != nil {
		logrus.Errorln("Error reading cursor file:", err)
		return cursorTime, false, err
	}

	var lastLine string
	for _, row := range bytes.Split(fileContent, []byte("\n")) {
		line := string(row)
		if line == "" {
			continue
		}
		lastLine = line
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

	fileContentRow, err := io.ReadAll(postFileHandler)
	if err != nil {
		logrus.Errorln("Error reading file:", err)
		return nil, 0, err
	}

	fileContent := bytes.Split(fileContentRow, []byte("\n"))

	nonEmptyContent := make([][]byte, 0)
	for _, line := range fileContent {
		if len(line) > 0 {
			nonEmptyContent = append(nonEmptyContent, line)
		}
	}
	fileContent = nonEmptyContent
	lineCount := len(fileContent)

	return fileContent, lineCount, nil
}
