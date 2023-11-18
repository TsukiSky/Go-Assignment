package logger

import (
	"log"
	"os"
)

var Logger *log.Logger
var PerformanceLogger *log.Logger

func Init(homeworkId, logFilename, logPrefix string) {
	if Logger == nil {
		logFile, err := os.Create(homeworkId + "\\" + logFilename)
		if err != nil {
			log.Fatal("unable to create log file:", err)
		}
		Logger = log.New(logFile, logPrefix, log.LstdFlags)
	}
}

func InitPerformanceLog(homeworkId string) {
	if PerformanceLogger == nil {
		logFile, err := os.Create(homeworkId + "\\" + "performance.log")
		if err != nil {
			log.Fatal("unable to create performance log file:", err)
		}
		PerformanceLogger = log.New(logFile, "performance:", log.LstdFlags)
	}
}
