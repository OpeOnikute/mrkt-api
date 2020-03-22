package handlers

import (
	"log"
	"os"

	"github.com/google/logger"
)

const applogPath = "./logs/app.log"
const errlogPath = "./logs/error.log"

// AppLogger ...
var AppLogger *logger.Logger

// ErrorLogger ...
var ErrorLogger *logger.Logger

// InitLogger ...
func InitLogger() {

	appF := openLogFile(applogPath)
	errF := openLogFile(errlogPath)

	// Log to system log and a log file, Info logs don't write to stdout.
	AppLogger = logger.Init("AppLogger", false, true, appF)
	ErrorLogger = logger.Init("ErrorLogger", false, true, errF)

	logger.SetFlags(log.LstdFlags)
}

// CloseLoggers ...
func CloseLoggers() {
	AppLogger.Close()
	ErrorLogger.Close()
}

func openLogFile(logPath string) *os.File {
	lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		logger.Fatalf("Failed to open log file: %v", err)
	}
	return lf
}
