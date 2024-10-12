package model

const (
	SEPARATOR = byte('\t')
)

var (
	COMMAND_STORAGE_FOLDER      = ".malamtime/commands"
	COMMAND_PRE_STORAGE_FILE    = COMMAND_STORAGE_FOLDER + "/pre.txt"
	COMMAND_POST_STORAGE_FILE   = COMMAND_STORAGE_FOLDER + "/post.txt"
	COMMAND_CURSOR_STORAGE_FILE = COMMAND_STORAGE_FOLDER + "/cursor.txt"
)

func InitFolder(inTestMode bool) {
	if inTestMode {
		COMMAND_STORAGE_FOLDER = ".malamtime-testing/commands"
	}

	COMMAND_PRE_STORAGE_FILE = COMMAND_STORAGE_FOLDER + "/pre.txt"
	COMMAND_POST_STORAGE_FILE = COMMAND_STORAGE_FOLDER + "/post.txt"
	COMMAND_CURSOR_STORAGE_FILE = COMMAND_STORAGE_FOLDER + "/cursor.txt"
}

type GinGraphQLContextType struct {
	IP     string
	UserID int
}
