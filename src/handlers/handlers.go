package handlers

import (
	misc "Misc"
	"data"
	"encoding/json"
	"net/http"
)

//CreateUser registers a user
func CreateUser(w http.ResponseWriter, request *http.Request) {

	user := &data.User{}

	err := json.NewDecoder(request.Body).Decode(user)

	if err != nil {
		misc.Warning("Can not parse request", err)
		w.WriteHeader(http.StatusForbidden)
		Respond(w, Message(false, "Cannot parse request"))
		return
	}

	respond, err := user.Register()

	if err != nil {
		misc.Warning(err)
		w.WriteHeader(http.StatusForbidden)
		Respond(w, Message(false, "Can not register"))
		return
	}
	w.WriteHeader(http.StatusCreated)
	Respond(w, respond)
}
