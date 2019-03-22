package data

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
)

//Token represents the token struct
type Token struct {
	UserId    uint32
	FirstName string
	jwt.StandardClaims
}

//Verify represents the incoming token
type Verify struct {
	Token string `json:"token"`
}

//User represents the User struct
type User struct {
	ID        uint32    `json:"id"`
	UUID      string    `json:"-"`
	FirstName string    `json:"firstname"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}

//Users has logged users
var Users []string

//Register registers a user to the system
func (user *User) Register() (response map[string]interface{}, err error) {
	check := user.IsUser()
	if check {
		err = errors.New("User already exists")
		return
	}
	query := "insert into users (uuid, email, password, created_at, firstname) values ($1, $2, $3, $4, $5) returning id, uuid, created_at, firstname"
	stmt, err := Db.Prepare(query)

	if err != nil {
		return
	}

	defer stmt.Close()

	err = stmt.QueryRow(createUUID(), user.Email, Encrypt(user.Password), time.Now(), user.FirstName).Scan(&user.ID, &user.UUID, &user.CreatedAt, &user.FirstName)

	response = map[string]interface{}{"id": user.ID, "uuid": user.UUID, "createdAt": user.CreatedAt, "email": user.Email, "firstname": user.FirstName}
	return
}

//IsUser checks whether a use is already registered or not.
func (user *User) IsUser() (available bool) {
	var num int
	Db.QueryRow("select COUNT(*) from users where email=$1 limit 1", user.Email).Scan(&num)

	if num == 0 {
		available = false
	} else {
		available = true
	}
	return
}

// LoginUser tries to login and returns if it's valid
func (user *User) LoginUser(email string, password string) (valid bool, err error) {
	err = Db.QueryRow("select id, password, uuid, created_at, firstname from users where email=$1", email).Scan(&user.ID, &user.Password, &user.UUID, &user.CreatedAt, &user.FirstName)
	if err != nil {
		err = errors.New("Email Not Found")
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		err = errors.New("Invalid credentials")
		return
	}

	valid = true
	return
}
