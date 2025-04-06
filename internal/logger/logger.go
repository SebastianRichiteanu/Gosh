package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// Log levels
const (
	LevelDebug = "DEBUG"
	LevelInfo  = "INFO"
	LevelWarn  = "WARN"
	LevelError = "ERROR"
)

var order = map[string]int{
	LevelDebug: 0,
	LevelInfo:  1,
	LevelWarn:  2,
	LevelError: 3,
}

type Logger struct {
	logFile     *os.File
	logFilePath string
	logLevel    string
}

func NewLogger(logFilePath, logLevel string) (*Logger, error) {
	logger := Logger{
		logFilePath: logFilePath,
		logLevel:    logLevel,
	}

	if err := logger.setLoggingFile(); err != nil {
		return nil, err
	}

	return &logger, nil
}

func (l *Logger) setLoggingFile() error {
	var err error
	l.logFile, err = os.OpenFile(l.logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	log.SetOutput(l.logFile)

	return nil
}

func (l *Logger) ensureLoggingFile() {
	_, err := os.Stat(l.logFilePath)
	if err == nil {
		return
	}

	// If the file was deleted, close the old reference
	l.logFile.Close()
	l.logFile = nil
	// And set it again
	if err := l.setLoggingFile(); err != nil {
		panic(err)
	}
}

// shouldLog checks if a message should be logged based on the current level
func (l *Logger) shouldLog(level string) bool {
	return order[level] >= order[l.logLevel]
}

// formatFields processes key-value pairs into "key=value"
func formatFields(fields []any) string {
	if len(fields) == 0 {
		return ""
	}

	if len(fields)%2 != 0 {
		return " | Invalid log fields: Must be key-value pairs"
	}

	var parts []string
	for i := 0; i < len(fields); i += 2 {
		key, ok := fields[i].(string)
		if !ok {
			return " | Invalid log key: Keys must be strings"
		}
		parts = append(parts, fmt.Sprintf("%s=%v", key, fields[i+1]))
	}

	return " | " + strings.Join(parts, "; ")
}

// logMessage formats and writes log messages
func (l *Logger) logMessage(level, message string, fields []any) {
	if !l.shouldLog(level) {
		return
	}

	l.ensureLoggingFile()

	log.Printf("[%s] %s%s", level, message, formatFields(fields))
}

// Info logs an info message
func (l *Logger) Debug(msg string, fields ...any) {
	l.logMessage(LevelDebug, msg, fields)
}

// Info logs an info message
func (l *Logger) Info(msg string, fields ...any) {
	l.logMessage(LevelInfo, msg, fields)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields ...any) {
	l.logMessage(LevelWarn, msg, fields)
}

// Error logs an error message
func (l *Logger) Error(msg string, fields ...any) {
	l.logMessage(LevelError, msg, fields)
}

// Close closes the log file
func (l *Logger) Close() {
	if l.logFile != nil {
		l.logFile.Close()
	}
}
