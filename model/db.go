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
	DB = db
	logrus.Traceln("DB connected")
}

func Clean() {
	DB.Close()
}
