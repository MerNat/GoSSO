package data

import (
	"errors"
	"time"
	"github.com/dgrijalva/jwt-go"
)

//Token represents the token struct
type Token struct {
	Uuid string
	jwt.StandardClaims
}

//User represents the User struct
type User struct {
	ID        uint32    `json:"id"`
	UUID      string    `json:"-"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}

//Register registers a user to the system
func (user *User) Register() (response map[string]interface{}, err error) {
	check := user.IsUser()

	if check {
		err = errors.New("User already exists")
		return
	}
	query := "insert into users (uuid, email, password, created_at) values ($1, $2, $3, $4) returning id, uuid, created_at"
	stmt, err := Db.Prepare(query)

	if err != nil {
		return
	}

	defer stmt.Close()

	err = stmt.QueryRow(createUUID(), user.Email, Encrypt(user.Password), time.Now()).Scan(&user.ID, &user.UUID, &user.CreatedAt)

	response = map[string]interface{}{"id":user.ID,"uuid":user.UUID,"createdAt":user.CreatedAt,"email": user.Email}
	return
}

//IsUser checks whether a use is already registered or not.
func (user *User) IsUser() (available bool) {
	return
}
