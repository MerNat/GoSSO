package handlers

import (
	"encoding/json"
	"net/http"
)

// Message returns the message about to send to the client
func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

// Respond returns data to the client
func Respond(w http.ResponseWriter, data map[string]interface{}) {
	json.NewEncoder(w).Encode(data)
}
