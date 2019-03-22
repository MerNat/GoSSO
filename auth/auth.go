package auth

import (
	"net/http"
)

// JwtAuth is a global JWT-Authenticator
var JwtAuth = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
		return
	})
}
