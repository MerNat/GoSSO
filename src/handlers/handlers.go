package handlers

import (
	"encoding/json"
	"net/http"

	misc "github.com/MerNat/GoSSO/src/Misc"
	"github.com/MerNat/GoSSO/src/data"
	"github.com/dgrijalva/jwt-go"
)

//CreateUser registers a user
func CreateUser(w http.ResponseWriter, request *http.Request) {

	user := &data.User{}

	err := json.NewDecoder(request.Body).Decode(user)

	if err != nil {
		misc.Warning("Can not parse request", err)
		w.WriteHeader(http.StatusBadRequest)
		Respond(w, Message(false, "Cannot parse request"))
		return
	}

	respond, err := user.Register()

	if err != nil {
		misc.Warning(err)
		w.WriteHeader(http.StatusBadRequest)
		Respond(w, Message(false, err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	Respond(w, respond)
}

//Login logs and generate a token
func Login(w http.ResponseWriter, request *http.Request) {
	user := &data.User{}
	err := json.NewDecoder(request.Body).Decode(user)

	if err != nil {
		response := Message(false, "Cant parse incomming data")
		w.WriteHeader(http.StatusBadRequest)
		Respond(w, response)
		return
	}

	_, err = user.LoginUser(user.Email, user.Password)

	if err != nil {
		response := Message(false, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		Respond(w, response)
		return
	}
	tk := &data.Token{UserId: user.ID}
	tk.ExpiresAt = misc.Config.JwtExpires
	tk.Issuer = misc.Config.JwtIssuer
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tk)
	tokenString, _ := token.SignedString([]byte(misc.Config.JwtSecret))
	//before appending check if user exists in map
	data.Users = append(data.Users, tokenString)
	w.WriteHeader(http.StatusOK)
	Respond(w, map[string]interface{}{"token": tokenString})
}
