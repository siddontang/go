package log

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

//log level, from low to high, more high means more serious
const (
	LevelTrace = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

const (
	Ltime  = 1 << iota //time format "2006/01/02 15:04:05"
	Lfile              //file.go:123
	Llevel             //[Trace|Debug|Info...]
)

var LevelName [6]string = [6]string{"Trace", "Debug", "Info", "Warn", "Error", "Fatal"}

const TimeFormat = "2006/01/02 15:04:05"

const maxBufPoolSize = 16

type Logger struct {
	sync.Mutex

	level int
	flag  int

	handler Handler

	quit chan struct{}
	msg  chan []byte

	bufs [][]byte

	wg sync.WaitGroup

	closed bool
}

//new a logger with specified handler and flag
func New(handler Handler, flag int) *Logger {
	var l = new(Logger)

	l.level = LevelInfo
	l.handler = handler

	l.flag = flag

	l.quit = make(chan struct{})
	l.closed = false

	l.msg = make(chan []byte, 1024)

	l.bufs = make([][]byte, 0, 16)

	l.wg.Add(1)
	go l.run()

	return l
}

//new a default logger with specified handler and flag: Ltime|Lfile|Llevel
func NewDefault(handler Handler) *Logger {
	return New(handler, Ltime|Lfile|Llevel)
}

func newStdHandler() *StreamHandler {
	h, _ := NewStreamHandler(os.Stdout)
	return h
}

var std = NewDefault(newStdHandler())

func (l *Logger) run() {
	defer l.wg.Done()
	for {
		select {
		case msg := <-l.msg:
			l.handler.Write(msg)
			l.putBuf(msg)
		case <-l.quit:
			//we must log all msg
			if len(l.msg) == 0 {
				return
			}
		}
	}
}

func (l *Logger) popBuf() []byte {
	l.Lock()
	var buf []byte
	if len(l.bufs) == 0 {
		buf = make([]byte, 0, 1024)
	} else {
		buf = l.bufs[len(l.bufs)-1]
		l.bufs = l.bufs[0 : len(l.bufs)-1]
	}
	l.Unlock()

	return buf
}

func (l *Logger) putBuf(buf []byte) {
	l.Lock()
	if len(l.bufs) < maxBufPoolSize {
		buf = buf[0:0]
		l.bufs = append(l.bufs, buf)
	}
	l.Unlock()
}

func (l *Logger) Close() {
	if l.closed {
		return
	}
	l.closed = true

	close(l.quit)

	l.wg.Wait()

	l.quit = nil

	l.handler.Close()
}

//set log level, any log level less than it will not log
func (l *Logger) SetLevel(level int) {
	l.level = level
}

//a low interface, maybe you can use it for your special log format
//but it may be not exported later......
func (l *Logger) Output(callDepth int, level int, format string, v ...interface{}) {
	if l.closed {
		//closed
		return
	}

	if l.level > level {
		return
	}

	buf := l.popBuf()

	if l.flag&Ltime > 0 {
		now := time.Now().Format(TimeFormat)
		buf = append(buf, '[')
		buf = append(buf, now...)
		buf = append(buf, "] "...)
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

		buf = append(buf, file...)
		buf = append(buf, ':')

		buf = strconv.AppendInt(buf, int64(line), 10)
		buf = append(buf, ' ')
	}

	if l.flag&Llevel > 0 {
		buf = append(buf, '[')
		buf = append(buf, LevelName[level]...)
		buf = append(buf, "] "...)
	}

	s := fmt.Sprintf(format, v...)

	buf = append(buf, s...)

	if s[len(s)-1] != '\n' {
		buf = append(buf, '\n')
	}

	l.msg <- buf
}

//log with Trace level
func (l *Logger) Trace(format string, v ...interface{}) {
	l.Output(2, LevelTrace, format, v...)
}

//log with Debug level
func (l *Logger) Debug(format string, v ...interface{}) {
	l.Output(2, LevelDebug, format, v...)
}

//log with info level
func (l *Logger) Info(format string, v ...interface{}) {
	l.Output(2, LevelInfo, format, v...)
}

//log with warn level
func (l *Logger) Warn(format string, v ...interface{}) {
	l.Output(2, LevelWarn, format, v...)
}

//log with error level
func (l *Logger) Error(format string, v ...interface{}) {
	l.Output(2, LevelError, format, v...)
}

//log with fatal level
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
