package main

import (
	"flag"
	"github.com/gopusher/gateway/comet"
	"github.com/gopusher/gateway/log"
	"github.com/joho/godotenv"
)

func main() {
	filename := getArgs()

	log.Info("Load config file: %s", *filename)
	if err := godotenv.Load(*filename); err != nil {
		panic(err)
	}

	comet.Run()
}

func getArgs() (filename *string) {
	filename = flag.String("c", "./config.ini", "set config file path")
	flag.Parse()

	return
}
