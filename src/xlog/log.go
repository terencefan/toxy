package xlog

import (
	"context"
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

var DefaultLog LeveledLogger

// time, level, laddr, raddr, message
var defaultFormat = "%s %-8s %s - %s <%s> %s\n"

type LeveledLogger interface {
	Debug(ctx context.Context, p string)
	Info(ctx context.Context, p string)
	Warning(ctx context.Context, p string)
	Error(ctx context.Context, p string)
}

type Logger struct {
	level  int16
	format string
}

func get_default(
	ctx context.Context, key string, defaultVal interface{},
) (val interface{}) {
	if val = ctx.Value(key); val == nil {
		val = defaultVal
	}
	return val
}

func (self *Logger) Log(
	ctx context.Context,
	level int16,
	p string,
	args ...interface{},
) {
	if self.level > level {
		return
	}

	var (
		name  = get_default(ctx, "name", "N")
		laddr = get_default(ctx, "laddr", "N")
		raddr = get_default(ctx, "raddr", "N")
	)

	fmt.Printf(
		self.format,
		time.Now().Format("01-02 15:04:05.000"),
		level2str[level],
		laddr,
		raddr,
		name,
		fmt.Sprintf(p, args...),
	)
}

func (self *Logger) Debug(ctx context.Context, p string) {
	self.Log(ctx, DEBUG, p)
}

func (self *Logger) Info(ctx context.Context, p string) {
	self.Log(ctx, INFO, p)
}

func (self *Logger) Warning(ctx context.Context, p string) {
	self.Log(ctx, WARNING, p)
}

func (self *Logger) Error(ctx context.Context, p string) {
	self.Log(ctx, ERROR, p)
}

func MakeLogger(level int16) (log *Logger) {
	log = &Logger{}
	log.level = level
	log.format = defaultFormat
	return
}

func init() {
	DefaultLog = MakeLogger(DEBUG)
}

func Debug(ctx context.Context, p string) {
	DefaultLog.Debug(ctx, p)
}

func Info(ctx context.Context, p string) {
	DefaultLog.Info(ctx, p)
}

func Warning(ctx context.Context, p string) {
	DefaultLog.Warning(ctx, p)
}

func Error(ctx context.Context, p string) {
	DefaultLog.Error(ctx, p)
}
