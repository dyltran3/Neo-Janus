package logger

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	logFile     *os.File
	logMutex    sync.Mutex
	logWriter   *bufio.Writer
	logDirPath  = "./data/logs"
	logFilePath = filepath.Join(logDirPath, "core.log")
)

// Init initialize logging system
func Init() error {
	logMutex.Lock()
	defer logMutex.Unlock()

	// Create log directory if not exists
	if err := os.MkdirAll(logDirPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Open log file in append mode
	var err error
	logFile, err = os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	// Create buffered writer for better performance
	logWriter = bufio.NewWriter(logFile)
	return nil
}

// logMsg write log message with timestamp
func logMsg(level, format string, v ...interface{}) {
	logMutex.Lock()
	defer logMutex.Unlock()

	msg := fmt.Sprintf(format, v...)
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logLine := fmt.Sprintf("[%s] [%s] %s\n", timestamp, level, msg)

	// Write to console
	fmt.Print(logLine)

	// Write to file if initialized
	if logWriter != nil {
		if _, err := logWriter.WriteString(logLine); err != nil {
			fmt.Printf("[ERROR] Failed to write log: %v\n", err)
		}
	}
}

// Info log info level message
func Info(format string, v ...interface{}) {
	logMsg("INFO", format, v...)
}

// Error log error level message
func Error(format string, v ...interface{}) {
	logMsg("ERROR", format, v...)
}

// Debug log debug level message (development only)
func Debug(format string, v ...interface{}) {
	logMsg("DEBUG", format, v...)
}

// Close flush buffered content and close log file
func Close() error {
	logMutex.Lock()
	defer logMutex.Unlock()

	if logWriter != nil {
		if err := logWriter.Flush(); err != nil {
			fmt.Printf("[ERROR] Failed to flush log buffer: %v\n", err)
			return err
		}
	}

	if logFile != nil {
		if err := logFile.Close(); err != nil {
			fmt.Printf("[ERROR] Failed to close log file: %v\n", err)
			return err
		}
	}

	return nil
}
