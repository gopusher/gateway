package main

import (
	"github.com/joho/godotenv"
	"github.com/gopusher/gateway/monitor"
	"github.com/gopusher/gateway/comet"
	"flag"
	"github.com/gopusher/gateway/log"
)

func main() {
	filename, isMonitor := getArgs()

	log.Info("Load config file: %s", *filename)
	godotenv.Load(*filename)

	if *isMonitor {
		monitor.Run()

		return
	}

	comet.Run()
}

func getArgs() (filename *string, isMonitor *bool) {
	filename = flag.String("c", "./config.ini", "set config file path")
	//是否为 monitor 节点
	isMonitor = flag.Bool("m", false, "if running with monitor model")
	flag.Parse()

	return
}
