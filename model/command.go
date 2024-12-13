package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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
	Shell     string       `json:"shell"`
	SessionID int64        `json:"sid"`
	Command   string       `json:"cmd"`
	Main      string       `json:"main"`
	Hostname  string       `json:"hn"`
	Username  string       `json:"un"`
	Time      time.Time    `json:"t"`
	EndTime   time.Time    `json:"et"`
	Result    int          `json:"result"`
	Phase     CommandPhase `json:"phase"`

	// Only work in file
	RecordingTime time.Time `json:"-"`
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

func (c Command) DoSavePre() error {
	if err := ensureStorageFolder(); err != nil {
		return err
	}

	buf, err := c.ToLine(time.Now())
	if err != nil {
		return err
	}

	preFile := os.ExpandEnv(os.ExpandEnv("$HOME/" + COMMAND_PRE_STORAGE_FILE))
	f, err := os.OpenFile(preFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		err = fmt.Errorf("failed to open pre-command storage file: %v", err)
		logrus.Errorln(err)
		return err
	}
	defer f.Close()

	if _, err := f.Write(buf); err != nil {
		err = fmt.Errorf("failed to write to pre-command storage file: %v", err)
		logrus.Errorln(err)
		return err
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
	buf, err := c.ToLine(time.Now())
	if err != nil {
		return err
	}

	postFile := os.ExpandEnv("$HOME/" + COMMAND_POST_STORAGE_FILE)
	f, err := os.OpenFile(postFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open pre-command storage file: %v", err)
	}
	defer f.Close()

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

func (c Command) IsNil() bool {
	if c.Command != "" {
		return false
	}
	if c.SessionID != 0 {
		return false
	}
	if c.Username != "" {
		return false
	}
	if c.Shell != "" {
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

func (cmd Command) GetUniqueKey() string {
	return fmt.Sprintf("%s|%d|%s|%s", cmd.Shell, cmd.SessionID, cmd.Command, cmd.Username)
}

func (cmd Command) ToLine(recordingTime time.Time) (buf []byte, err error) {
	buf, err = json.Marshal(cmd)
	if err != nil {
		return
	}
	timestamp := recordingTime.UnixNano()
	timestampBytes := []byte(fmt.Sprintf("%d", timestamp))
	buf = append(buf, SEPARATOR)
	buf = append(buf, timestampBytes...)
	buf = append(buf, '\n')
	return
}

func (cmd *Command) FromLine(line string) (recordingTime time.Time, err error) {
	parts := strings.Split(line, string(SEPARATOR))
	if len(parts) != 2 {
		err = fmt.Errorf("Invalid line format in pre-command file: %s", line)
		logrus.Errorln(err)
		return
	}

	err = json.Unmarshal([]byte(parts[0]), cmd)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal command: %v, %s", err, parts[0])
		logrus.Errorln(err)
		return
	}

	unixNano, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		err = fmt.Errorf("failed to parse timestamp: %v, %s", err, parts[1])
		logrus.Errorln(err)
		return
	}
	recordingTime = time.Unix(0, unixNano)
	cmd.RecordingTime = recordingTime
	return
}

func (cmd *Command) FromLineBytes(line []byte) (recordingTime time.Time, err error) {
	parts := bytes.Split(line, []byte{SEPARATOR})
	if len(parts) != 2 {
		err = fmt.Errorf("Invalid line format in pre-command file: %s", line)
		logrus.Errorln(err)
		return
	}

	err = json.Unmarshal(parts[0], cmd)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal command: %v, %s", err, string(parts[0]))
		logrus.Errorln(err)
		return
	}

	unixNano, err := strconv.ParseInt(string(parts[1]), 10, 64)
	if err != nil {
		err = fmt.Errorf("failed to parse timestamp: %v, %s", err, string(parts[1]))
		logrus.Errorln(err)
		return
	}
	recordingTime = time.Unix(0, unixNano)
	cmd.RecordingTime = recordingTime
	return
}

func (cmd Command) FindClosestCommand(commandList []*Command, withSameKey bool) *Command {
	closestPreCommand := new(Command)
	minTimeDiff := int64(^uint64(0) >> 1) // Max int64 value

	myKey := cmd.GetUniqueKey()

	for _, c := range commandList {
		if withSameKey {
			ck := c.GetUniqueKey()
			if myKey != ck {
				continue
			}
		}

		timeDiff := cmd.Time.Unix() - c.Time.Unix()
		if timeDiff >= 0 && timeDiff < minTimeDiff {
			minTimeDiff = timeDiff
			closestPreCommand = c
		}
	}

	return closestPreCommand
}
