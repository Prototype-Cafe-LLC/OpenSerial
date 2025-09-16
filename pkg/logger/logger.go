package logger

import (
	"log"
	"os"
	"time"
)

// Logger provides simple logging functionality
type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
}

// NewLogger creates a new logger instance
func NewLogger() *Logger {
	return &Logger{
		infoLogger:  log.New(os.Stdout, "[INFO] ", log.LstdFlags),
		errorLogger: log.New(os.Stderr, "[ERROR] ", log.LstdFlags),
		debugLogger: log.New(os.Stdout, "[DEBUG] ", log.LstdFlags),
	}
}

// Info logs an info message
func (l *Logger) Info(format string, v ...interface{}) {
	l.infoLogger.Printf(format, v...)
}

// Error logs an error message
func (l *Logger) Error(format string, v ...interface{}) {
	l.errorLogger.Printf(format, v...)
}

// Debug logs a debug message
func (b *Logger) Debug(format string, v ...interface{}) {
	b.debugLogger.Printf(format, v...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(format string, v ...interface{}) {
	l.errorLogger.Fatalf(format, v...)
}

// LogConnection logs connection events
func (l *Logger) LogConnection(event, address string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	l.Info("[%s] %s: %s", timestamp, event, address)
}

// LogDataTransfer logs data transfer events
func (l *Logger) LogDataTransfer(direction string, bytes int) {
	l.Debug("Data %s: %d bytes", direction, bytes)
}

// LogError logs errors with context
func (l *Logger) LogError(context string, err error) {
	l.Error("%s: %v", context, err)
}

// Global logger instance
var Default = NewLogger()
