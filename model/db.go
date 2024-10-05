package model

import (
	"os"

	"github.com/nutsdb/nutsdb"
	"github.com/sirupsen/logrus"
)

var DB *nutsdb.DB

type GinGraphQLContextType struct {
	IP     string
	UserID int
}

func InitDB() {
	localDBPath := os.ExpandEnv("$HOME/.malamtime/db")
	db, err := nutsdb.Open(
		nutsdb.DefaultOptions,
		nutsdb.WithDir(localDBPath),
	)
	if err != nil {
		logrus.Panicln(err)
		return
	}
	// if config.GetRuntimeConfig().Debug {
	// 	entClient = entClient.Debug()
	// }

	err = db.Update(func(tx *nutsdb.Tx) error {
		if !tx.ExistBucket(nutsdb.DataStructureBTree, activeBucket) {
			err := tx.NewBucket(nutsdb.DataStructureBTree, activeBucket)
			if err != nil {
				return err
			}
		}
		if !tx.ExistBucket(nutsdb.DataStructureBTree, archivedBucket) {
			err := tx.NewBucket(nutsdb.DataStructureBTree, archivedBucket)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		logrus.Panicln(err)
		return
	}

	DB = db
	logrus.Traceln("DB connected")
}

func Clean() {
	DB.Close()
}
