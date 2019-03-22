package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

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
	tk.ExpiresAt = time.Now().Add(time.Hour * time.Duration(misc.Config.JwtExpires)).Unix()
	tk.Issuer = misc.Config.JwtIssuer
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tk)
	tokenString, _ := token.SignedString([]byte(misc.Config.JwtSecret))
	//before appending check if user exists in map
	data.Users = append(data.Users, tokenString)
	w.WriteHeader(http.StatusOK)
	Respond(w, map[string]interface{}{"token": tokenString})
}

// IsAuthorized middleware for verifying token
func IsAuthorized(w http.ResponseWriter, r *http.Request) {
	tk := data.Verify{}
	err := json.NewDecoder(r.Body).Decode(&tk)
	if err != nil {
		fmt.Print(err.Error())
		response := Message(false, "Cant parse incomming data")
		w.WriteHeader(http.StatusBadRequest)
		Respond(w, response)
		return
	}

	if tk.Token == "" {
		response := Message(false, "wrong data format")
		w.WriteHeader(http.StatusBadRequest)
		Respond(w, response)
		return
	}

	sort.Strings(data.Users)
	i := sort.SearchStrings(data.Users, tk.Token)
	if i >= len(data.Users) || i < 0 {
		response := Message(false, "User not logged in")
		w.WriteHeader(http.StatusBadRequest)
		Respond(w, response)
		return
	}

	token, err := jwt.Parse(tk.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte(misc.Config.JwtSecret), nil
	})

	if err != nil {
		fmt.Println(err.Error())
		response := Message(false, "Malformed auth token")
		w.WriteHeader(http.StatusBadRequest)
		Respond(w, response)
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		w.WriteHeader(http.StatusOK)
		Respond(w, map[string]interface{}{"data": claims["UserId"]})
		return
	}
}
