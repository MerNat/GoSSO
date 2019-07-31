package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	misc "github.com/MerNat/GoSSO/Misc"
	"github.com/MerNat/GoSSO/data"
	"github.com/dgrijalva/jwt-go"
)

var user data.User

//CreateUser registers a user
func CreateUser(w http.ResponseWriter, request *http.Request) {
	json.NewDecoder(request.Body).Decode(&user)
	_, err := user.Register()

	if err != nil {
		misc.Warning(err)
		w.WriteHeader(http.StatusBadRequest)
		Respond(w, Message(false, err.Error()))
		return
	}
	json.NewEncoder(w).Encode(&user)
}

//Login logs and generate a token
func Login(w http.ResponseWriter, request *http.Request) {
	err := json.NewDecoder(request.Body).Decode(&user)

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
	tk := &data.Token{UserId: user.ID, FirstName: user.FirstName}
	tk.ExpiresAt = time.Now().Add(time.Hour * time.Duration(misc.Config.JwtExpires)).Unix()
	tk.Issuer = misc.Config.JwtIssuer
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tk)
	tokenString, _ := token.SignedString([]byte(misc.Config.JwtSecret))
	//before appending check if user exists in slice
	sort.Strings(data.Users)
	i := sort.SearchStrings(data.Users, tokenString)
	if i > 0 && i < len(data.Users) {
		w.WriteHeader(http.StatusOK)
		Respond(w, map[string]interface{}{"token": data.Users[i]})
		return
	}
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
		fmt.Println(i)
		response := Message(false, "User not logged in")
		w.WriteHeader(http.StatusBadRequest)
		Respond(w, response)
		return
	}

	token, err := jwt.Parse(tk.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte(misc.Config.JwtSecret), nil
	})

	if err != nil {
		data.Users = append(data.Users[:i], data.Users[i+1:]...)
		response := Message(false, "Malformed auth token")
		w.WriteHeader(http.StatusBadRequest)
		Respond(w, response)
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		w.WriteHeader(http.StatusOK)
		Respond(w, map[string]interface{}{"userId": claims["UserId"], "userName": claims["FirstName"]})
		return
	}
}
