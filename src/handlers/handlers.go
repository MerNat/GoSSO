package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	misc "github.com/MerNat/GoSSO/src/Misc"
	"github.com/MerNat/GoSSO/src/data"
	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
)

var (
	user        = &data.User{}
	tokenString = token.SignedString([]byte(misc.Config.JwtSecret))
	tk          = &data.Token{UserId: user.ID}
	token       = jwt.NewWithClaims(jwt.SigningMethodHS256, tk)
)

//CreateUser registers a user
func CreateUser(w http.ResponseWriter, request *http.Request) {
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
	tk.ExpiresAt = misc.Config.JwtExpires
	tk.Issuer = misc.Config.JwtIssuer
	//before appending check if user exists in map
	data.Users = append(data.Users, tokenString)
	w.WriteHeader(http.StatusOK)
	Respond(w, map[string]interface{}{"token": tokenString})
}

// IsAuthorized middleware for verifying token
func IsAuthorized(ep func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["token"] != nil {
			token, err := jwt.Parse(r.Header["token"][0], func(tk *jwt.Token) (interface{}, error) {
				if _, ok := tk.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("something wrong happened")
				}
				return tokenString, nil
			})

			if err != nil {
				fmt.Fprintf(w, err.Error())
			}

			if token.Valid {
				var user data.User
				mapstructure.Decode(token.Claims, &user)
				vars := mux.Vars(r)
				email := vars["email"]
				if email != user.Email {
					Respond(w, map[string]interface{}{"Error": "Invalid authorization token - Does not match user email"})
					return
				}
				ep(w, r)
			}
		} else {
			fmt.Fprintf(w, "Not Authorized")
		}
	})
}
