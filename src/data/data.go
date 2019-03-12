package data

import (
	misc "Misc"
	"crypto/rand"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

//Db has a connection to db
var Db *sql.DB

func init() {
	var err error
	varSettings := "dbname=" + misc.Config.DbName + " user=" + misc.Config.DbUser + " port=" + misc.Config.DbPort +
		" host=" + misc.Config.DbAddress + " password=" + misc.Config.DbPassword
	Db, err = sql.Open("postgres", varSettings)

	if err != nil {
		misc.Error("Can not connect to DB", err)
	}
}

//Creates a 36 Character UUID
func createUUID() (uuid string) {
	u := new([16]byte)
	_, err := rand.Read(u[:])
	if err != nil {
		misc.Error("Cannot generate UUID", err)
	}
	u[8] = (u[8] | 0x40) & 0x7F
	u[6] = (u[6] & 0xF) | (0x4 << 4)
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	return
}

// Encrypt encypts a string with sha1 algorithm
func Encrypt(plaintext string) (cryptedtext string) {
	// cryptedtext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
	cryptedtext = string(hashPassword)
	return
}
