package log

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

const (
	LevelTrace = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

const (
	Ltime  = iota << 1 //time format "2006/01/02 15:04:05"
	Lfile              //file.go:123
	Llevel             //[Trace|Debug|Info...]
)

var LevelName [6]string = [6]string{"Trace", "Debug", "Info", "Warn", "Error", "Fatal"}

const TimeFormat = "2006/01/02 15:04:05"

type Logger struct {
	level int
	flag  int

	handler Handler

	quit chan struct{}
	msg  chan []byte
}

func New(handler Handler, flag int) *Logger {
	var l = new(Logger)

	l.level = LevelInfo
	l.handler = handler

	l.flag = flag

	l.quit = make(chan struct{})

	l.msg = make(chan []byte, 1024)

	go l.run()

	return l
}

func NewDefault(handler Handler) *Logger {
	return New(handler, Ltime|Lfile|Llevel)
}

func newStdHandler() *StreamHandler {
	h, _ := NewStreamHandler(os.Stdout)
	return h
}

var std = NewDefault(newStdHandler())

func (l *Logger) run() {
	for {
		select {
		case msg := <-l.msg:
			l.handler.Write(msg)
		case <-l.quit:
			l.handler.Close()
		}
	}
}

func (l *Logger) Close() {
	if l.quit == nil {
		return
	}

	close(l.quit)
	l.quit = nil
}

func (l *Logger) SetLevel(level int) {
	l.level = level
}

func (l *Logger) Output(callDepth int, level int, format string, v ...interface{}) {
	if l.level > level {
		return
	}

	buf := make([]byte, 0, 1024)

	if l.flag&Ltime > 0 {
		now := time.Now().Format(TimeFormat)
		buf = append(buf, now...)
		buf = append(buf, " "...)
	}

	if l.flag&Lfile > 0 {
		_, file, line, ok := runtime.Caller(callDepth)
		if !ok {
			file = "???"
			line = 0
		} else {
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					file = file[i+1:]
					break
				}
			}
		}

		buf = append(buf, fmt.Sprintf("%s:%d ", file, line)...)
	}

	if l.flag&Llevel > 0 {
		buf = append(buf, fmt.Sprintf("[%s] ", LevelName[level])...)
	}

	s := fmt.Sprintf(format, v...)

	buf = append(buf, s...)

	if s[len(s)-1] != '\n' {
		buf = append(buf, "\n"...)
	}

	l.msg <- buf
}

func (l *Logger) Trace(format string, v ...interface{}) {
	l.Output(2, LevelTrace, format, v...)
}

func (l *Logger) Debug(format string, v ...interface{}) {
	l.Output(2, LevelDebug, format, v...)
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.Output(2, LevelInfo, format, v...)
}

func (l *Logger) Warn(format string, v ...interface{}) {
	l.Output(2, LevelWarn, format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.Output(2, LevelError, format, v...)
}

func (l *Logger) Fatal(format string, v ...interface{}) {
	l.Output(2, LevelFatal, format, v...)
}

func SetLevel(level int) {
	std.SetLevel(level)
}

func Trace(format string, v ...interface{}) {
	std.Output(2, LevelTrace, format, v...)
}

func Debug(format string, v ...interface{}) {
	std.Output(2, LevelDebug, format, v...)
}

func Info(format string, v ...interface{}) {
	std.Output(2, LevelInfo, format, v...)
}

func Warn(format string, v ...interface{}) {
	std.Output(2, LevelWarn, format, v...)
}

func Error(format string, v ...interface{}) {
	std.Output(2, LevelError, format, v...)
}

func Fatal(format string, v ...interface{}) {
	std.Output(2, LevelFatal, format, v...)
}
