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

var DefaultLog Logger

type Logger interface {
	Debug(p string, args ...interface{})
	Info(p string, args ...interface{})
	Warning(p string, args ...interface{})
	Error(p string, args ...interface{})
}

type ConsoleLog struct {
	level int
}

func (self *ConsoleLog) Log(level, p string, args ...interface{}) {
	prefix := fmt.Sprintf("[%s][%s]", time.Now().Format("01-02 15:04:05.000"), level)
	fmt.Printf(prefix+" "+p+"\n", args...)
}

func (self *ConsoleLog) Debug(p string, args ...interface{}) {
	if self.level > DEBUG {
		return
	}
	self.Log("DEBUG", p, args...)
}

func (self *ConsoleLog) Info(p string, args ...interface{}) {
	if self.level > INFO {
		return
	}
	self.Log("INFO", p, args...)
}

func (self *ConsoleLog) Warning(p string, args ...interface{}) {
	if self.level > WARNING {
		return
	}
	self.Log("WARNING", p, args...)
}

func (self *ConsoleLog) Error(p string, args ...interface{}) {
	if self.level > ERROR {
		return
	}
	self.Log("ERROR", p, args...)
}

func MakeConsoleLog(level int) (log *ConsoleLog) {
	log = &ConsoleLog{}
	log.level = level
	return
}

func init() {
	DefaultLog = MakeConsoleLog(DEBUG)
	DefaultLog.Info("xlog has been initialized.")
}

func Debug(p string, args ...interface{}) {
	DefaultLog.Debug(p, args...)
}

func Info(p string, args ...interface{}) {
	DefaultLog.Info(p, args...)
}

func Warning(p string, args ...interface{}) {
	DefaultLog.Warning(p, args...)
}

func Error(p string, args ...interface{}) {
	DefaultLog.Error(p, args...)
}
