package model

import (
	"encoding/json"
	"errors"
	"fmt"
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
	Shell        string
	SessionID    int64
	Command      string
	Main         string
	Hostname     string
	Username     string
	Time         time.Time
	EndTime      time.Time
	Result       int
	Phase        CommandPhase
	SentToServer bool
}

func (c Command) DoSavePre() error {
	return DB.Update(func(tx *nutsdb.Tx) error {
		key := c.getDBKey(true)
		buf, err := json.Marshal(c)
		if err != nil {
			return err
		}
		return tx.Put(activeBucket, []byte(key), buf, nutsdb.Persistent)
	})
}

func (c Command) DoUpdate(result int) error {
	return DB.Update(func(tx *nutsdb.Tx) error {
		keys, vals, err := tx.GetAll(activeBucket)
		if err != nil {
			return err
		}

		var matchedKey []byte

		for i, key := range keys {
			var item Command
			err := json.Unmarshal(vals[i], &item)
			if err != nil {
				return err
			}
			if c.IsPairPreCommand(item) {
				matchedKey = key
				break
			}
		}

		if matchedKey == nil {
			return errors.New("pre command could not found")
		}

		c.Result = result
		c.EndTime = time.Now()
		buf, err := json.Marshal(c)
		if err != nil {
			return err
		}
		if err := tx.Put(archivedBucket, matchedKey, buf, nutsdb.Persistent); err != nil {
			return err
		}
		return tx.Delete(activeBucket, matchedKey)
	})
}

func (c Command) IsPairPreCommand(target Command) bool {
	if !c.IsSame(target) {
		return false
	}
	if c.Phase != CommandPhasePre {
		return false
	}
	// one command not possible to run for 10 days, right?
	if c.Time.Before(time.Now().Add(time.Hour * 24 * 10)) {
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

func GetArchivedList(length int) (keys [][]byte, values []Command, err error) {
	valueBytes := make([][]byte, length)

	err = DB.View(func(tx *nutsdb.Tx) error {
		keys, valueBytes, err = tx.GetAll(archivedBucket)
		return err
	})

	if err != nil {
		return
	}

	for i, vb := range valueBytes {
		var cmd Command
		if err = json.Unmarshal(vb, &cmd); err != nil {
			return nil, nil, err
		}
		values[i] = cmd
	}

	return keys, values, err
}

func GetArchievedCount() (count int, err error) {
	err = DB.View(func(tx *nutsdb.Tx) error {
		size, err := tx.LSize(archivedBucket, []byte("*"))

		if err != nil {
			return err
		}
		count = size
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
