package model

import (
	"context"
	"fmt"
	"os"

	"github.com/malamtime/cli/ent"
	"github.com/sirupsen/logrus"

	_ "github.com/mattn/go-sqlite3"
)

var EntClient *ent.Client

type GinGraphQLContextType struct {
	IP     string
	UserID int
}

func InitDB() {
	localDBPath := os.ExpandEnv("$HOME/.malamtime/local.db")
	entClient, err := ent.Open("sqlite3", fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", localDBPath))
	if err != nil {
		logrus.Panicln(err)
		return
	}
	// if config.GetRuntimeConfig().Debug {
	// 	entClient = entClient.Debug()
	// }

	err = entClient.
		Schema.
		Create(
			context.Background(),
		)
	if err != nil {
		logrus.Panicln(err)
	}

	EntClient = entClient
	logrus.Infoln("ent client connected")
}

func Clean() {
	EntClient.Close()
}
