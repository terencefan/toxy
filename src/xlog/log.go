package xlog

import (
	"fmt"
	"time"
)

const (
	DEBUG = (iota - 1) * 10
	INFO
	WARNING
	ERROR
)

var level2str = map[int16]string{
	DEBUG:   "DEBUG",
	INFO:    "INFO",
	WARNING: "WARNING",
	ERROR:   "ERROR",
}

var str2level = map[string]int16{
	"DEBUG":   DEBUG,
	"INFO":    INFO,
	"WARNING": WARNING,
	"ERROR":   ERROR,
}

var DefaultLog LeveledLogger

// time, level, laddr, raddr, message
var defaultFormat = "%s [%s] %s\n"

type LeveledLogger interface {
	Debug(p string)
	Info(p string)
	Warning(p string)
	Error(p string)
	Level(level int16)
}

type Logger struct {
	level  int16
	format string
}

func (self *Logger) Log(
	level int16,
	p string,
	args ...interface{},
) {
	if self.level > level {
		return
	}

	fmt.Printf(
		self.format,
		time.Now().Format("2006-01-02 15:04:05.000"),
		level2str[level],
		fmt.Sprintf(p, args...),
	)
}

func (self *Logger) Level(level int16) {
	self.level = level
}

func (self *Logger) Debug(p string) {
	self.Log(DEBUG, p)
}

func (self *Logger) Info(p string) {
	self.Log(INFO, p)
}

func (self *Logger) Warning(p string) {
	self.Log(WARNING, p)
}

func (self *Logger) Error(p string) {
	self.Log(ERROR, p)
}

func MakeLogger(level int16) (log *Logger) {
	log = &Logger{}
	log.level = level
	log.format = defaultFormat
	return
}

func init() {
	DefaultLog = MakeLogger(INFO)
}

func Level(level int16) {
	DefaultLog.Level(level)
}

func LevelString(level string) {
	if l, ok := str2level[level]; !ok {
		Warning("invalid log level: " + level)
	} else {
		Level(l)
	}
}

func Debug(p string) {
	DefaultLog.Debug(p)
}

func Info(p string) {
	DefaultLog.Info(p)
}

func Warning(p string) {
	DefaultLog.Warning(p)
}

func Error(p string) {
	DefaultLog.Error(p)
}
