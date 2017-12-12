package main

import (
	"os"
	"runtime"
	"ting"
	"toxy"
	"xlog"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	config     = kingpin.Flag("config", "Config file.").Short('c').Default("toxy.ini").String()
	fast_reply = kingpin.Flag("fast-reply", "Fast reply ping request").Bool()
	level      = kingpin.Flag("level", "Log level").Short('l').Default("INFO").String()
	ping       = kingpin.Flag("ping", "Send ping request to all the defined backends").Bool()
)

func runPing(conf *toxy.Config) {
	if err := ting.Run(conf); err != nil {
		xlog.Error(err.Error())
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

func run(conf *toxy.Config) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var toxy = toxy.NewToxy(conf)
	if *fast_reply {
		toxy.FastReply()
	}
	toxy.Serve()
}

func main() {
	kingpin.Parse()

	xlog.LevelString(*level)

	conf, err := toxy.LoadConfig(*config)
	if err != nil {
		panic(err)
	}

	if *ping {
		runPing(conf)
	} else {
		run(conf)
	}
}
