package xmetric

import (
	"net"
	"time"
)

type buffered_statsd struct {
	conn    net.Conn
	addr    string
	prefix  string
	timeout time.Duration

	flush_period    time.Duration
	idle_period     time.Duration
	max_buffer_size int
	max_queue_size  int
	queue           chan string
	buf             []byte

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
	if len(self.buf) == 0 {
		return
	}

	defer func() {
		self.buf = self.buf[:0]
	}()

	var err = self.connect()
	if err != nil {
		error_handler(err)
		return
	}

	_, err = self.conn.Write(self.buf)
	if err != nil {
		error_handler(err)
		return
	}
}

func (self *buffered_statsd) append(message string) {
	// NOTE ignore messages which are too long.
	if len(message) > self.max_buffer_size {
		return
	}
	if len(self.buf)+len(message) > self.max_buffer_size {
		self.flush()
	}
	self.buf = append(self.buf, message...)
}

func (self *buffered_statsd) connect() error {
	if self.conn != nil {
		return nil
	}
	conn, err := net.DialTimeout("tcp", self.addr, self.timeout)
	self.conn = conn
	return err
}

func (self *buffered_statsd) open() {
	if self.opened {
		return
	}
	self.opened = true

	go func() {
		// flush buf periodicly
		var ticker = time.NewTicker(self.flush_period)
		// close conn and stop goroutine
		var countdown = time.NewTimer(self.idle_period)

		defer ticker.Stop()
		defer countdown.Stop()

		for {
			select {
			case message := <-self.queue:
				countdown.Reset(self.idle_period)
				self.append(message)
			case <-ticker.C:
				self.flush()
			case <-countdown.C:
				if self.conn != nil {
					self.conn.Close()
					self.conn = nil
				}
				self.opened = false
				return
			}
		}
	}()
}

func (self *buffered_statsd) close() {
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
		buf:             make([]byte, 0, conf.MaxBufferSize),
	}
	return h
}

func InitBufferedStatsd(options ...Option) {
	handler = NewBufferedStatsd(options...)
}
