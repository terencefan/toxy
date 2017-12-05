package xlog

type Handler interface {
	Log(level int16, p string, args ...interface{})
}
