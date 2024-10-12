package model

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/nutsdb/nutsdb"
)

// pre commands here
const activeBucket = "active"

// post commands here
// will delete after sync to server
const archivedBucket = "archived"

type CommandPhase int

const (
	CommandPhasePre  = 0
	CommandPhasePost = 1
)

type Command struct {
	Shell     string
	SessionID int64
	Command   string
	Main      string
	Hostname  string
	Username  string
	Time      time.Time
	EndTime   time.Time
	Result    int
	Phase     CommandPhase
}

func ensureStorageFolder() error {
	storageFolder := os.ExpandEnv("$HOME/" + COMMAND_STORAGE_FOLDER)
	if _, err := os.Stat(storageFolder); os.IsNotExist(err) {
		if err := os.MkdirAll(storageFolder, 0755); err != nil {
			return fmt.Errorf("failed to create command storage folder: %v", err)
		}
	}

	return nil
}

func addTimestampToCommandBuf(buf []byte) []byte {
	timestamp := time.Now().UnixNano()
	timestampBytes := []byte(fmt.Sprintf("%d", timestamp))
	buf = append(buf, SEPARATOR)
	buf = append(buf, timestampBytes...)
	buf = append(buf, '\n')
	return buf
}

func (c Command) DoSavePre() error {
	if err := ensureStorageFolder(); err != nil {
		return err
	}

	buf, err := json.Marshal(c)
	if err != nil {
		return err
	}

	preFile := os.ExpandEnv("$HOME/" + COMMAND_PRE_STORAGE_FILE)
	f, err := os.OpenFile(preFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open pre-command storage file: %v", err)
	}
	defer f.Close()
	buf = addTimestampToCommandBuf(buf)

	if _, err := f.Write(buf); err != nil {
		return fmt.Errorf("failed to write to pre-command storage file: %v", err)
	}
	return nil
}

func (c Command) DoUpdate(result int) error {
	if err := ensureStorageFolder(); err != nil {
		return err
	}

	c.Phase = CommandPhasePost
	c.Result = result
	c.EndTime = time.Now()
	buf, err := json.Marshal(c)
	if err != nil {
		return err
	}

	postFile := os.ExpandEnv("$HOME/" + COMMAND_POST_STORAGE_FILE)
	f, err := os.OpenFile(postFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open pre-command storage file: %v", err)
	}
	defer f.Close()
	buf = addTimestampToCommandBuf(buf)

	if _, err := f.Write(buf); err != nil {
		return fmt.Errorf("failed to write to pre-command storage file: %v", err)
	}
	return nil
}

func (c Command) IsPairPreCommand(target Command) bool {
	if !c.IsSame(target) {
		return false
	}
	if c.Phase != CommandPhasePre {
		return false
	}
	// one command not possible to run for 10 days, right?
	if c.Time.Before(time.Now().Add(-time.Hour * 24 * 10)) {
		return false
	}
	return true
}

func (c Command) IsSame(target Command) bool {
	if c.Shell != target.Shell {
		return false
	}
	if c.Command != target.Command {
		return false
	}
	if c.SessionID != target.SessionID {
		return false
	}
	if c.Username != target.Username {
		return false
	}
	return true
}

func (c Command) getDBKey(withUUid bool) string {
	key := fmt.Sprintf("%s:%d", c.Shell, c.SessionID)
	if withUUid {
		key += ":" + uuid.New().String()
	}
	return key
}

func GetArchivedList(keys [][]byte) (values []Command, err error) {
	valueBytes := make([][]byte, len(keys))

	err = DB.View(func(tx *nutsdb.Tx) error {
		valueBytes, err = tx.MGet(archivedBucket, keys...)
		return err
	})

	if err != nil {
		return
	}

	for _, vb := range valueBytes {
		var cmd Command
		if err = json.Unmarshal(vb, &cmd); err != nil {
			return nil, err
		}
		values = append(values, cmd)
	}

	return values, err
}

func GetArchievedCount() (keys [][]byte, err error) {
	err = DB.View(func(tx *nutsdb.Tx) error {
		bucketKeys, err := tx.GetKeys(archivedBucket)

		if err != nil {
			return err
		}
		keys = bucketKeys
		return nil
	})
	return
}

func CleanArchievedData(keys [][]byte) error {
	return DB.Update(func(tx *nutsdb.Tx) error {
		for _, key := range keys {
			if err := tx.Delete(archivedBucket, key); err != nil {
				return err
			}
		}
		return nil
	})
}
