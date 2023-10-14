package logger

import (
	"log"
	"os"
)

var Logger *log.Logger

func init() {
	if Logger == nil {
		logFile, err := os.Create("log.log")
		if err != nil {
			log.Fatal("unable to create log file:", err)
		}
		Logger = log.New(logFile, "assignment 1: ", log.LstdFlags)
	}
}
