package xmetric

import (
	"bytes"
	"net"
	"time"
)

type buffered_statsd struct {
	addr    string
	prefix  string
	timeout time.Duration

	flush_period    time.Duration
	idle_period     time.Duration
	max_buffer_size int
	max_queue_size  int
	queue           chan string
	buf             *bytes.Buffer

	opened bool
}

func (self *buffered_statsd) send(message string) error {
	self.open()
	self.queue <- message
	return nil
}

func (self *buffered_statsd) Timing(key string, delta int) error {
	var message = message(self.prefix, key, "timers", "%d|ms", delta)
	return self.send(message)
}

func (self *buffered_statsd) Count(key string, delta int) error {
	var message = message(self.prefix, key, "counters", "%d|c", delta)
	return self.send(message)
}

func (self *buffered_statsd) Gauge(key string, value int) error {
	var message = message(self.prefix, key, "gauge", "%d|g", value)
	return self.send(message)
}

func (self *buffered_statsd) flush() {
	if self.buf.Len() == 0 {
		return
	}

	defer func() {
		self.buf.Reset()
	}()

	conn, err := self.connect()
	if err != nil {
		error_handler(err)
		return
	}

	if _, err := self.buf.WriteTo(conn); err != nil {
		error_handler(err)
		return
	}
}

func (self *buffered_statsd) append(message string) {
	message += "\n"
	// NOTE ignore message which are too long.
	if len(message) > self.max_buffer_size {
		return
	}
	if self.buf.Len()+len(message) > self.max_buffer_size {
		self.flush()
	}
	self.buf.WriteString(message)
}

func (self *buffered_statsd) connect() (net.Conn, error) {
	return net.DialTimeout("udp", self.addr, self.timeout)
}

func (self *buffered_statsd) open() {
	if self.opened {
		return
	}
	self.opened = true

	go func() {
		// flush buf periodicly
		var ticker = time.NewTicker(self.flush_period)

		defer ticker.Stop()

		for {
			select {
			case message := <-self.queue:
				self.append(message)
			case <-ticker.C:
				self.flush()
			}
		}
	}()
}

func NewBufferedStatsd(options ...Option) (h *buffered_statsd) {
	var conf = NewConfig(options...)
	h = &buffered_statsd{
		addr:            conf.Addr,
		prefix:          conf.Prefix,
		timeout:         conf.Timeout,
		flush_period:    conf.FlushPeriod,
		idle_period:     time.Second * 5,
		max_buffer_size: conf.MaxBufferSize,
		max_queue_size:  conf.MaxQueueSize,
		queue:           make(chan string, conf.MaxQueueSize),
		buf:             bytes.NewBuffer([]byte{}),
	}
	return h
}

func InitBufferedStatsd(options ...Option) {
	handler = NewBufferedStatsd(options...)
}
