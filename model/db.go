package model

import (
	"bufio"
	"bytes"
	"context"
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

func GetPreCommandsTree(ctx context.Context) (result preCommandTree, err error) {
	ctx, span := modelTracer.Start(ctx, "db.getPreCmdsTree")
	defer span.End()
	preFilePath := os.ExpandEnv(fmt.Sprintf("%s/%s", "$HOME", COMMAND_PRE_STORAGE_FILE))
	preFileHandler, err := os.Open(preFilePath)
	if err != nil {
		logrus.Errorln("Failed to open pre-command file:", err)
		return nil, err
	}
	defer preFileHandler.Close()

	result = make(preCommandTree)
	scanner := bufio.NewScanner(preFileHandler)
	for scanner.Scan() {
		line := scanner.Bytes()
		cmd := new(Command)
		_, err := cmd.FromLineBytes(line)
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

	if err := scanner.Err(); err != nil {
		return result, err
	}

	return result, nil
}

func GetPreCommands(ctx context.Context) ([]*Command, error) {
	ctx, span := modelTracer.Start(ctx, "db.getPreCmds")
	defer span.End()
	preFilePath := os.ExpandEnv(fmt.Sprintf("%s/%s", "$HOME", COMMAND_PRE_STORAGE_FILE))
	preFileHandler, err := os.Open(preFilePath)
	if err != nil {
		logrus.Errorln("Failed to open pre-command file:", err)
		return nil, err
	}
	defer preFileHandler.Close()

	result := make([]*Command, 0)
	scanner := bufio.NewScanner(preFileHandler)

	for scanner.Scan() {
		raw := scanner.Bytes()
		if len(raw) == 0 {
			continue
		}
		cmd := new(Command)
		_, err := cmd.FromLineBytes(raw)
		if err != nil {
			logrus.Errorln("Invalid line parse in pre-command file:", string(raw), err)
			continue
		}
		result = append(result, cmd)
	}

	if err := scanner.Err(); err != nil {
		logrus.Errorln("Error reading file:", err)
		return nil, err
	}

	return result, nil
}

func GetLastCursor(ctx context.Context) (cursorTime time.Time, noCursorExist bool, err error) {
	ctx, span := modelTracer.Start(ctx, "db.getLastCursor")
	defer span.End()
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

func GetPostCommands(ctx context.Context) ([][]byte, int, error) {
	ctx, span := modelTracer.Start(ctx, "db.getPostCmds")
	defer span.End()
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
