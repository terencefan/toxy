package main

import (
	"fmt"
	"sync"

	. "github.com/stdrickforce/thriftgo/protocol"
)

type TestFunc func(proto Protocol) error

var wg sync.WaitGroup

func call(name string, args ...interface{}) TestFunc {
	return func(proto Protocol) error {
		return nil
	}
}

func processor(id int, fn TestFunc) {
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
		go processor(i, call("ping"))
		wg.Add(1)
	}
	wg.Wait()
}
