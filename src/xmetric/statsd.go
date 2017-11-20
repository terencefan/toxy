package xmetric

import (
	"fmt"
	"net"
	"time"
)

// TODO logger support
type statsd struct {
	conn    net.Conn
	addr    string
	prefix  string
	timeout time.Duration
}

func (self statsd) Timing(key string, delta int) error {
	var message = message(self.prefix, key, "timers", "%d|ms", delta)
	return self.send(message)
}

func (self statsd) Count(key string, delta int) error {
	var message = message(self.prefix, key, "counters", "%d|c", delta)
	return self.send(message)
}

func (self statsd) Gauge(key string, value int) error {
	var message = message(self.prefix, key, "gauge", "%d|g", value)
	return self.send(message)
}

func (self *statsd) connect() error {
	if self.conn != nil {
		return nil
	}
	var conn, err = net.DialTimeout(
		"udp",
		self.addr,
		self.timeout,
	)
	if err != nil {
		return err
	}
	self.conn = conn
	return nil
}

func message(prefix, stat, cate, format string, val int) (m string) {
	if len(prefix) > 0 {
		m = fmt.Sprintf(
			"%s.%s.%s:%s",
			prefix,
			cate,
			stat,
			fmt.Sprintf(format, val),
		)
	} else {
		m = fmt.Sprintf(
			"%s.%s:%s",
			cate,
			stat,
			fmt.Sprintf(format, val),
		)
	}
	return
}

func (self *statsd) send(message string) (e error) {
	_, e = fmt.Fprintf(self.conn, message)
	return
}

func NewStatsd(options ...Option) (h *statsd) {
	var conf = NewConfig(options...)
	h = &statsd{
		addr:    conf.Addr,
		prefix:  conf.Prefix,
		timeout: conf.Timeout,
	}
	h.connect()
	return
}

func InitStatsd(options ...Option) {
	handler = NewStatsd(options...)
}
