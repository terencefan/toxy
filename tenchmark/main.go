package main

import (
	"fmt"
	"sync"

	. "github.com/stdrickforce/thriftgo/protocol"
)

var wg sync.WaitGroup

func call(proto Protocol, name string, args ...interface{}) error {
	return nil
}

func processor(id int, fn func(proto Protocol) error) {
	defer wg.Done()

	// TODO build a real protocol
	var proto Protocol

	for i := 0; i < 20; i++ {
		fn(proto)
		fmt.Println(id, i)
	}
}

func main() {
	for i := 0; i < 10; i++ {
		fn := func(proto Protocol) error {
			return call(proto, "ping")
		}
		go processor(i, fn)
		wg.Add(1)
	}
	wg.Wait()
}
