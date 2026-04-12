package logger

import (
	"log"
	"os"
)

type Logger interface {
	Info(msg string)
	Error(msg string)
	Debug(msg string)
}

type defaultLogger struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	debugLog *log.Logger
}

func New() Logger {
	return &defaultLogger{
		infoLog:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime),
		errorLog: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		debugLog: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime),
	}
}

func (l *defaultLogger) Info(msg string) {
	l.infoLog.Println(msg)
}

func (l *defaultLogger) Error(msg string) {
	l.errorLog.Println(msg)
}

func (l *defaultLogger) Debug(msg string) {
	l.debugLog.Println(msg)
}
