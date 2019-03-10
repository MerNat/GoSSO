package data

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

//Token represents the token struct
type Token struct {
	UUID string
	jwt.StandardClaims
}

//User represents the User struct
type User struct {
	ID        uint32 `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Token     string `json:"token"`
	CreatedAt time.Time
}
