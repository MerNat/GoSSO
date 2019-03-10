package main

import (
	"encoding/json"
	"log"
	"os"
)

//Configuration holds the struct to the config
type Configuration struct {
	Address string
	Port    string
}

var outputLogger *log.Logger

//Config holds the values used to start the server
var Config Configuration

func init() {
	file, err := os.OpenFile("sso.log", os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalln("Can Access sso.log file", err)
	}

	outputLogger = log.New(file, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)
}

func load() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalln("Can not open config file", err)
	}
	err = json.NewDecoder(file).Decode(&Config)

	if err != nil {
		log.Fatalln("Can not decode config file", err)
	}
}

func danger(args ...interface{}) {
	outputLogger.SetPrefix("DANGER ")
	outputLogger.Println(args...)
}

func warning(args ...interface{}) {
	outputLogger.SetPrefix("WARNING ")
	outputLogger.Println(args...)
}

func error(args ...interface{}) {
	outputLogger.SetPrefix("ERROR ")
	outputLogger.Println(args...)
}

func info(args ...interface{}) {
	outputLogger.Println(args...)
}