package log

import (
	"fmt"
	"log"
	"os"
)

const (
	LevelTrace = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

var LevelName [6]string = [6]string{"Trace", "Debug", "Info", "Warn", "Error", "Fatal"}

type Logger struct {
	logger  *log.Logger
	level   int
	handler Handler
}

func New(handler Handler, flag int) *Logger {
	var l = new(Logger)
	l.logger = log.New(handler, "", flag) //log.LstdFlags|log.Lshortfile)
	l.level = LevelInfo
	l.handler = handler

	return l
}

func NewDefault(handler Handler) *Logger {
	return New(handler, log.LstdFlags|log.Lshortfile)
}

func newStdHandler() *StreamHandler {
	h, _ := NewDefaultStreamHandler(os.Stdout)
	return h
}

var std = NewDefault(newStdHandler())

func (l *Logger) Close() {
	l.handler.Close()
}

func (l *Logger) SetLevel(level int) {
	l.level = level
}

func (l *Logger) Output(callDepth int, level int, format string, v ...interface{}) {
	if l.level <= level {
		f := fmt.Sprintf("[%s] %s", LevelName[level], format)
		s := fmt.Sprintf(f, v...)
		l.logger.Output(callDepth, s)
	}
}

func (l *Logger) Write(s string) {
	l.logger.Output(3, s)
}

func (l *Logger) Trace(format string, v ...interface{}) {
	l.Output(3, LevelTrace, format, v...)
}

func (l *Logger) Debug(format string, v ...interface{}) {
	l.Output(3, LevelDebug, format, v...)
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.Output(3, LevelInfo, format, v...)
}

func (l *Logger) Warn(format string, v ...interface{}) {
	l.Output(3, LevelWarn, format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.Output(3, LevelError, format, v...)
}

func (l *Logger) Fatal(format string, v ...interface{}) {
	l.Output(3, LevelFatal, format, v...)
}

func SetLevel(level int) {
	std.SetLevel(level)
}

func Trace(format string, v ...interface{}) {
	std.Output(3, LevelTrace, format, v...)
}

func Debug(format string, v ...interface{}) {
	std.Output(3, LevelDebug, format, v...)
}

func Info(format string, v ...interface{}) {
	std.Output(3, LevelInfo, format, v...)
}

func Warn(format string, v ...interface{}) {
	std.Output(3, LevelWarn, format, v...)
}

func Error(format string, v ...interface{}) {
	std.Output(3, LevelError, format, v...)
}

func Fatal(format string, v ...interface{}) {
	std.Output(3, LevelFatal, format, v...)
}
