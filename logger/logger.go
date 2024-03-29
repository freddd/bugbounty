package logger

import (
	"github.com/fatih/color"
	"os"
)

const (
	info    = "INFO: "
	debug   = "DEBUG: "
	warn    = "WARN: "
	error   = "ERROR: "
	fatal   = "FATAL: "
	newline = "\n"
)

type Logger interface {
	Debug(format string, args ...interface{})
	Error(format string, args ...interface{})
	Fatal(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
}

type defaultLogger struct {
}

var DefaultLogger Logger
var blue = color.New(color.FgBlue)
var red = color.New(color.FgRed)
var yellow = color.New(color.FgYellow)
var green = color.New(color.FgGreen)

func init() {
	DefaultLogger = setupLogger()
}

func setupLogger() Logger {
	return &defaultLogger{}
}

func (dl *defaultLogger) Debug(msg string, args ...interface{}) {
	blue.Fprintf(os.Stdout, appendNewline(msg), args...)
}

func (dl *defaultLogger) Info(msg string, args ...interface{}) {
	green.Fprintf(os.Stdout, appendNewline(msg), args...)
}

func (dl *defaultLogger) Error(msg string, args ...interface{}) {
	red.Fprintf(os.Stderr, appendNewline(msg), args...)
}

func (dl *defaultLogger) Fatal(msg string, args ...interface{}) {
	red.Fprintf(os.Stderr, appendNewline(msg), args...)
	os.Exit(1)
}

func (dl *defaultLogger) Warn(msg string, args ...interface{}) {
	yellow.Fprintf(os.Stdout, appendNewline(msg), args...)
}

func appendNewline(msg string) string {
	return msg + newline
}
