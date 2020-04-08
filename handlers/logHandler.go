package handlers

import (
	"log"
	"os"

	"github.com/google/logger"
)

const applogPath = "./logs/app.log"
const errlogPath = "./logs/error.log"

// Logger ...
type Logger struct {
	AppLogger   *logger.Logger
	ErrorLogger *logger.Logger
}

func (l *Logger) Info(txt string) {

	f := l.openLogFile(applogPath)

	// Log to system log and a log file, Info logs don't write to stdout.
	appLogger := logger.Init("appLogger", false, true, f)
	defer appLogger.Close()

	logger.SetFlags(log.LstdFlags)

	appLogger.Info(txt)
}

// TODO: this should accept an actual error
func (l *Logger) Error(txt string) {

	f := l.openLogFile(applogPath)

	// Log to system log and a log file, Info logs don't write to stdout.
	errorLogger := logger.Init("errorLogger", false, true, f)
	defer errorLogger.Close()

	logger.SetFlags(log.LstdFlags)

	errorLogger.Error(txt)
}

func (l *Logger) openLogFile(logPath string) *os.File {
	lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		logger.Fatalf("Failed to open log file: %v", err)
		panic(err)
	}
	return lf
}
