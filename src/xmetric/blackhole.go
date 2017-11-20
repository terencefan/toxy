package xmetric

type blackhole struct {
}

func (self *blackhole) Timing(key string, delta int) error {
	return nil
}

func (self *blackhole) Count(key string, delta int) error {
	return nil
}

func (self *blackhole) Gauge(key string, value int) error {
	return nil
}
