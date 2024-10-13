package model

import (
	"fmt"
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

type GinGraphQLContextType struct {
	IP     string
	UserID int
}
