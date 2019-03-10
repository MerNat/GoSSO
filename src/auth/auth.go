package auth

import (
	"handlers"
	"net/http"
	"strings"
)

// JwtAuth is a global JWT-Authenticator
var JwtAuth = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rejectURL := []string{"/sso/register", "/sso/login"}
		currentPath := r.URL.Path

		for _, value := range rejectURL {
			if value == currentPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		response := make(map[string]interface{})
		tokenHeader := r.Header.Get("Authorization")

		if tokenHeader == "" {
			response = handlers.Message(false, "Missing Token")
			w.WriteHeader(http.StatusForbidden)
			handlers.Respond(w, response)
		}

		splitted := strings.Split(tokenHeader, " ")

		if len(splitted) != 2 {
			response = handlers.Message(false, "Invalid/Malformed Token")
			w.WriteHeader(http.StatusForbidden)
			handlers.Respond(w, response)
			return
		}

		theToken := splitted[1]
		

	})
}
