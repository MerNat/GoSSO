package auth

import (
	"context"
	"net/http"
	"strings"

	misc "github.com/MerNat/GoSSO/src/Misc"
	"github.com/MerNat/GoSSO/src/data"
	"github.com/MerNat/GoSSO/src/handlers"

	"github.com/dgrijalva/jwt-go"
)

// JwtAuth is a global JWT-Authenticator
var JwtAuth = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rejectURL := []string{"/sso/register", "/sso/login"}
		currentPath := r.URL.Path
		w.Header().Add("Content-Type", "application/json")
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

		tk := &data.Token{}

		token, err := jwt.ParseWithClaims(theToken, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(misc.Config.JwtSecret), nil
		})

		if err != nil {
			response = handlers.Message(false, "Malformed auth token")
			w.WriteHeader(http.StatusBadRequest)
			handlers.Respond(w, response)
			return
		}

		if !token.Valid {
			response = handlers.Message(false, "Token is not valid or expired")
			w.WriteHeader(http.StatusForbidden)
			handlers.Respond(w, response)
			return
		}

		ctx := context.WithValue(r.Context(), "user", tk.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
