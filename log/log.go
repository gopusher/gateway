package log

import (
	"log"
	"os"
)

var logger = log.New(os.Stdout, "", log.LstdFlags)

//Info log info data
func Info(format string, v ...interface{}) {
	logger.Printf("[Info] "+format, v...)
}

//Warn log warn data
func Warn(format string, v ...interface{}) {
	logger.Printf("\033[33m[Warn] "+format+"\033[0m", v...)
}

//Error log error info
func Error(format string, v ...interface{}) {
	logger.Printf("\033[31m[Error] "+format+"\033[0m", v...)
}
