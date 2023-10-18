package logger

import (
	"log"
	"os"
)

var Logger *log.Logger

func Init(logFilename, logPrefix string) {
	if Logger == nil {
		logFile, err := os.Create("hw1\\" + logFilename)
		if err != nil {
			log.Fatal("unable to create log file:", err)
		}
		Logger = log.New(logFile, logPrefix, log.LstdFlags)
	}
}
