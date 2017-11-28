package main

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	. "github.com/stdrickforce/thriftgo/protocol"
	. "github.com/stdrickforce/thriftgo/thrift"
	. "github.com/stdrickforce/thriftgo/transport"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type TestFunc func(proto Protocol) error

var wg sync.WaitGroup

var count int32 = 0

func call(name string, args ...interface{}) TestFunc {
	var writeMessageBody = func(proto Protocol) (err error) {
		if err = proto.WriteStructBegin("whatever"); err != nil {
			return
		}
		for i, arg := range args {
			index := int16(i + 1)
			switch v := arg.(type) {
			case int16:
				err = proto.WriteFieldBegin("i16", T_I16, index)
				err = proto.WriteI16(v)
			case int32:
				err = proto.WriteFieldBegin("i32", T_I32, index)
				err = proto.WriteI32(v)
			case int64:
				err = proto.WriteFieldBegin("i64", T_I64, index)
				err = proto.WriteI64(v)
			case string:
				err = proto.WriteFieldBegin("string", T_STRING, index)
				err = proto.WriteString(v)
			default:
				err = errors.New("unsupport type")
			}
			if err != nil {
				return
			}
		}
		if err = proto.WriteFieldStop(); err != nil {
			return
		}
		if err = proto.WriteStructEnd(); err != nil {
			return
		}
		if err = proto.WriteMessageEnd(); err != nil {
			return
		}
		if err = proto.Flush(); err != nil {
			return
		}
		return
	}

	return func(proto Protocol) (err error) {
		atomic.AddInt32(&count, 1)
		if count%100 == 0 {
			fmt.Println(count)
		}

		if err = proto.WriteMessageBegin(name, T_CALL, 0); err != nil {
			return
		}

		if err = writeMessageBody(proto); err != nil {
			return
		}
		if _, _, _, err = proto.ReadMessageBegin(); err != nil {
			return
		}
		if err = proto.Skip(T_STRUCT); err != nil {
			return
		}
		if err = proto.ReadMessageEnd(); err != nil {
			return
		}
		return
	}
}

func processor(id int, fn TestFunc, requests int) {
	defer wg.Done()
	defer func() {
		fmt.Println(id, "exit")
	}()

	// TODO build a real protocol
	var (
		trans Transport
		proto Protocol
	)

	tf := NewTSocketFactory("192.168.100.116:32327")
	// tw := NewTFramedTransportFactory(false, true)
	tw := TTransportWrapper
	pf := NewTBinaryProtocolFactory(true, true)

	trans = tf.GetTransport()
	trans = tw.GetTransport(trans)
	proto = pf.GetProtocol(trans)
	if err := trans.Open(); err != nil {
		panic(err)
	}
	defer trans.Close()

	for i := 0; i < requests; i++ {
		if err := fn(proto); err != nil {
			fmt.Println(id, err)
			return
		}
	}
}

var (
	concurrency = kingpin.Flag("concurrency", "Number of multiple requests to make at a time").Short('c').Default("1").Int()
	requests    = kingpin.Flag("requests", "Number of requests to perform").Short('n').Default("1").Int()
)

func main() {
	kingpin.Parse()
	fmt.Println(*concurrency, *requests)
	for i := 0; i < *concurrency; i++ {
		go processor(i, call("RevenueOrder:get", int16(2), int32(102), int64(201711100250000008)), *requests)
		wg.Add(1)
	}
	wg.Wait()
}
