package main

import (
	"log"
	"net/http"
	"time"

	misc "github.com/MerNat/GoSSO/Misc"
	"github.com/MerNat/GoSSO/auth"
	"github.com/MerNat/GoSSO/handlers"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/sso/register", handlers.CreateUser).Methods("POST")
	router.HandleFunc("/sso/login", handlers.Login).Methods("POST")
	router.HandleFunc("/sso/verify", handlers.IsAuthorized).Methods("POST")

	router.Use(auth.JwtAuth)

	server := &http.Server{
		Addr:           misc.Config.ServerAddress,
		Handler:        router,
		ReadTimeout:    time.Duration(misc.Config.ServerReadTimeout * int64(time.Second)),
		WriteTimeout:   time.Duration(misc.Config.ServerWriteTimeout * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	misc.Info("Server Started: ", misc.Config.ServerAddress)
	log.Fatal(server.ListenAndServe())
}
