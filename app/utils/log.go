package utils

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/sirupsen/logrus"
)

// Logger creates and returns a configured logrus.Logger instance.
// It creates a logs directory in the current working directory and generates a log file named with the current date.
// The log file is opened in append mode, and the log level is set to Debug.
// The timestamp format for logs is "2006-01-02 15:04:05".
func Logger() *logrus.Logger {
	// Get current time
	now := time.Now()
	// Set log file path
	logFilePath := ""
	if dir, err := os.Getwd(); err == nil {
		logFilePath = dir + "/logs/"
	}
	// Create logs directory
	if err := os.MkdirAll(logFilePath, 0777); err != nil {
		fmt.Println(err.Error())
	}
	// Set log file name with format "2006-01-02.log"
	logFileName := now.Format("2006-01-02") + ".log"

	// Combine full log file path
	fileName := path.Join(logFilePath, logFileName)
	// Check if log file exists, create if it doesn't
	if _, err := os.Stat(fileName); err != nil {
		if _, err := os.Create(fileName); err != nil {
			fmt.Println(err.Error())
		}
	}

	// Open log file in append mode
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}

	// Create a new logrus.Logger instance
	logger := logrus.New()
	// Set log output to file
	logger.Out = src
	// Set log level to Debug
	logger.SetLevel(logrus.DebugLevel)
	// Set log format
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})

	// Return the configured logger instance
	return logger
}
