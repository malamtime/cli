package commands

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

var loggerFile *os.File

func SetupLogger() {
	// TODO: change to size based logger selector
	logFilePath := os.ExpandEnv("$HOME/.malamtime/log.log")

	f, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		os.Stdout.WriteString(err.Error())
		fmt.Errorf("[MalamTime] on setup logger error:%s \n", err.Error())
	}
	loggerFile = f
	logrus.SetReportCaller(true)
	logrus.SetOutput(loggerFile)
	logrus.Infoln("Setting up logger with version: ", commitID)
}

func CloseLogger() {
	loggerFile.Close()
}