package log

import (
	"log"
	"os"
)

var logger = log.New(os.Stdout, "", log.LstdFlags)

func Info(format string, v ...interface{}) {
	logger.Printf("[Info] "+format, v...)
}

func Warn(format string, v ...interface{}) {
	logger.Printf("\033[33m[Warn] "+format+"\033[0m", v...)
}

func Error(format string, v ...interface{}) {
	logger.Printf("\033[31m[Error] "+format+"\033[0m", v...)
}
