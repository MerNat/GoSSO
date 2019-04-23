package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

type UserTest struct {
	Firstname string `json:"firstname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type Users struct {
	Users []UserTest `json:"users"`
}

var mux2 *http.ServeMux
var writer *httptest.ResponseRecorder
var users Users
var apiURL string

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	os.Exit(code)
}

func setUp() {
	apiURL = "http://0.0.0.0:8080"
	file, err := os.Open("tests.json")
	if err != nil {
		log.Fatalln("Can not open config file", err)
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&users)

	if err != nil {
		log.Fatalln("Can not decode config file", err)
	}
	mux2 = http.NewServeMux()
	writer = httptest.NewRecorder()
}

func TestHandlePost(t *testing.T) {
	endpoint := "/sso/register"
	u, _ := url.ParseRequestURI(apiURL)
	u.Path = endpoint
	urlPath := u.String()
	for i := 0; i < len(users.Users); i++ {
		json, _ := json.Marshal(users.Users[i])
		request, _ := http.Post(urlPath, "application/json", bytes.NewBuffer(json))
		if request.StatusCode != 201 {
			t.Errorf("Response code is %v\n", request.StatusCode)
			break
		}
	}
}
