package misc

import (
	"encoding/json"
	"log"
	"os"
)

//Configuration holds the struct to the config
type Configuration struct {
	ServerAddress string
	ServerPort    string
	DbAddress     string
	DbPort        string
	DbUser        string
	DbName        string
	DbPassword    string
}

var outputLogger *log.Logger

//Config holds the config used to start the whole server
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

//Danger outputs and prints a danger errLog.
func Danger(args ...interface{}) {
	outputLogger.SetPrefix("DANGER ")
	outputLogger.Println(args...)
}

func Warning(args ...interface{}) {
	outputLogger.SetPrefix("WARNING ")
	outputLogger.Println(args...)
}

func Error(args ...interface{}) {
	outputLogger.SetPrefix("ERROR ")
	outputLogger.Fatalln(args...)
}

func Info(args ...interface{}) {
	outputLogger.Println(args...)
}
