package common

import (
	"log"
	"os"
)

//Log var for log
var Log *log.Logger

//NewLog create log File
func NewLog(logpath string) {
	println("LogFile: " + logpath)
	file, err := os.Create(logpath)
	if err != nil {
		panic(err)
	}
	Log = log.New(file, "", log.LstdFlags|log.Lshortfile)
}
