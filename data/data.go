package data

import (
	"crypto/rand"
	// "database/sql"
	"github.com/jinzhu/gorm"
	"fmt"

	misc "github.com/MerNat/GoSSO/Misc"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"os"
)

//Db has a connection to db
var Db *gorm.DB

func init() {
	var err error
	Db, err = gorm.Open(
		"postgres", 
		os.ExpandEnv("host=${HOST} user=${USER} dbname=${DBNAME} sslmode=disable password=${PASSWORD}"))

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
