package xmetric

import (
	"fmt"
	"math/rand"
	"time"
)

type Handler interface {
	Timing(key string, delta int) error
	Count(key string, delta int) error
	Gauge(key string, value int) error
}

var (
	handler       Handler
	error_handler func(e error)
)

func init() {
	handler = &blackhole{}
	error_handler = func(e error) {
		fmt.Println(e)
	}
}

func check(rate float32) bool {
	if rate >= 1.0 {
		return true
	} else if rate <= 0 {
		return false
	}
	r := rand.New(rand.NewSource(time.Now().Unix()))
	return r.Float32() < rate
}

// TODO custom mixer
func mix(category, key string) string {
	return fmt.Sprintf("%s.%s", category, key)
}

func Init(dsn string) {
	// TODO parse dsn and initialize handler.
}

func InitWithOptions(dsn string, options map[string]string) {
	// TODO parse dsn and initialize handler with options.
}

func Timing(category, key string, delta int) (e error) {
	key = mix(category, key)
	return handler.Timing(key, delta)
}

func Count(category, key string, delta int) (e error) {
	key = mix(category, key)
	return handler.Count(key, delta)
}

func Incr(category, key string) (e error) {
	key = mix(category, key)
	return handler.Count(key, 1)
}

func Decr(category, key string) (e error) {
	key = mix(category, key)
	return handler.Count(key, -1)
}

func Gauge(category, key string, value int) (e error) {
	key = mix(category, key)
	return handler.Gauge(key, value)
}

func TimingWithSampling(category, key string, delta int, rate float32) error {
	if check(rate) {
		return Timing(category, key, delta)
	} else {
		return nil
	}
}

func GaugeWithSampling(category, key string, value int, rate float32) error {
	if check(rate) {
		return Gauge(category, key, value)
	} else {
		return nil
	}
}

func CountWithSampling(category, key string, delta int, rate float32) error {
	if check(rate) {
		return Count(category, key, delta)
	} else {
		return nil
	}
}

func IncrWithSampling(category, key string, rate float32) error {
	if check(rate) {
		return Incr(category, key)
	} else {
		return nil
	}
}

func DecrWithSampling(category, key string, rate float32) error {
	if check(rate) {
		return Decr(category, key)
	} else {
		return nil
	}
}
