package xmetric

import (
	"fmt"
	"time"
)

type config struct {
	Addr          string
	Prefix        string
	Timeout       time.Duration
	FlushPeriod   time.Duration
	MaxBufferSize int
	MaxQueueSize  int
}

type Option func(*config)

func NewConfig(options ...Option) (c *config) {
	c = &config{
		Addr:          ":8125",
		Prefix:        "",
		Timeout:       time.Second * 5,
		FlushPeriod:   time.Second,
		MaxBufferSize: 1024,
		MaxQueueSize:  16,
	}
	for _, option := range options {
		option(c)
	}
	return
}

func Address(host string, port int) Option {
	return Option(func(c *config) {
		c.Addr = fmt.Sprintf("%s:%d", host, port)
	})
}

func Prefix(prefix string) Option {
	return Option(func(c *config) {
		c.Prefix = prefix
	})
}

func Timeout(timeout time.Duration) Option {
	return Option(func(c *config) {
		c.Timeout = timeout
	})
}

func FlushPeriod(period time.Duration) Option {
	return Option(func(c *config) {
		c.FlushPeriod = period
	})
}

func MaxBufferSize(size int) Option {
	return Option(func(c *config) {
		c.MaxBufferSize = size
	})
}

func MaxQueueSize(size int) Option {
	return Option(func(c *config) {
		c.MaxQueueSize = size
	})
}
