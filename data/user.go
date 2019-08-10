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
	FirstName string    `gorm:"column:firstname"`
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
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return
	}

	stmt := Db.Table("users").Create(
		&User{
			UUID:      createUUID(),
			FirstName: user.FirstName,
			Email:     user.Email,
			Password:  string(passwordHash)})

	if stmt.Error != nil {
		err = stmt.Error
		return
	}

	defer stmt.Close()

	response = map[string]interface{}{"uuid": user.UUID, "createdAt": user.CreatedAt, "email": user.Email, "firstname": user.FirstName}
	return response, nil
}

//IsUser checks whether a use is already registered or not.
func (user *User) IsUser() (available bool) {
	var num int
	err := Db.Table("users").Where("email = ?", user.Email).Count(&num).Error

	if err != nil {
		return false
	}

	if num == 0 {
		available = false
	} else {
		available = true
	}
	return
}

// LoginUser tries to login and returns if it's valid
func (user *User) LoginUser(email string, password string) (valid bool, err error) {
	err = Db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return false, err
	}
	// fmt.Println(password, user.Password)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return false, err
	}

	valid = true
	return
}
