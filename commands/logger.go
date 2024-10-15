package commands

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

var loggerFile *os.File

var SKIP_LOGGER_SETTINGS = false

func SetupLogger(baseFolder string) {
	if SKIP_LOGGER_SETTINGS {
		logrus.SetReportCaller(true)
		return
	}
	// TODO: change to size based logger selector
	logFilePath := baseFolder + "/log.log"

	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		err := os.MkdirAll(baseFolder, 0755)
		if err != nil {
			fmt.Errorf("[MalamTime] failed to create log directory: %s\n", err.Error())
			return
		}

		_, err = os.Create(logFilePath)
		if err != nil {
			fmt.Errorf("[MalamTime] failed to create log file: %s\n", err.Error())
			return
		}
	}

	f, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		os.Stdout.WriteString(err.Error())
		fmt.Errorf("[MalamTime] on setup logger error:%s \n", err.Error())
	}
	loggerFile = f
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetOutput(loggerFile)
	logrus.Traceln("Setting up logger with version: ", commitID)
}

func CloseLogger() {
	if SKIP_LOGGER_SETTINGS {
		return
	}
	logrus.Traceln("going to close...")
	loggerFile.Close()
}
