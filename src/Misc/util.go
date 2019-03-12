package misc

import (
	"encoding/json"
	"log"
	"os"

	"github.com/dgrijalva/jwt-go"
)

//Configuration holds the struct to the config
type Configuration struct {
	ServerAddress      string
	ServerReadTimeout  int64
	ServerWriteTimeout int64
	DbAddress          string
	DbPort             string
	DbUser             string
	DbName             string
	DbPassword         string
	JwtSecret          string
	JwtExpires         int64
	JwtIssuer          string
}

var outputLogger *log.Logger

//Config holds the config used to start the whole server
var Config Configuration

// StandardClaim holds JwtClaim
var StandardClaim jwt.StandardClaims

func init() {
	// file, err := os.OpenFile("sso.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	// if err != nil {
	// 	log.Fatalln("Can Access sso.log file", err)
	// }
	load()
	outputLogger = log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)
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
	StandardClaim = jwt.StandardClaims{
		ExpiresAt: Config.JwtExpires,
		Issuer:    Config.JwtIssuer,
	}
}

//Danger outputs and prints a danger errLog.
func Danger(args ...interface{}) {
	outputLogger.SetPrefix("DANGER ")
	outputLogger.Println(args...)
}

//Warning warns
func Warning(args ...interface{}) {
	outputLogger.SetPrefix("WARNING ")
	outputLogger.Println(args...)
}

//Error prints the message and terminate
func Error(args ...interface{}) {
	outputLogger.SetPrefix("ERROR ")
	outputLogger.Fatalln(args...)
}

//Info prints inf
func Info(args ...interface{}) {
	outputLogger.SetPrefix("INFO ")
	outputLogger.Println(args...)
}
