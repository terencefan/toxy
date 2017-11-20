package xtransport

type blackhole struct {
}

func (self *blackhole) Write(m []byte) (int, error) {
	return 0, nil
}

func (self *blackhole) Read(m []byte) (int, error) {
	return 0, nil
}

func (self *blackhole) Close() error {
	return nil
}

func (self *blackhole) Flush() error {
	return nil
}

var TBlackHole = &blackhole{}
