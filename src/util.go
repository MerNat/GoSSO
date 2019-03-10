package main

import (
	"log"
	"os"
)

var outputLogger *log.Logger

func init() {
	file, err := os.OpenFile("sso.log", os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalln("Can Access sso.log file", err)
	}

	outputLogger = log.New(file, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)
}
