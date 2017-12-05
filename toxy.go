package main

import (
	"runtime"
	"toxy"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	config     = kingpin.Flag("config", "Config file.").Short('c').Default("toxy.ini").String()
	fast_reply = kingpin.Flag("fast-reply", "Fast reply ping request").Bool()
)

func main() {
	kingpin.Parse()

	runtime.GOROOT()

	var toxy = toxy.NewToxy()
	if err := toxy.LoadConfig(*config); err != nil {
		panic(err)
	}
	if *fast_reply {
		toxy.FastReply()
	}
	toxy.Serve()
}
