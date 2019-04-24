package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
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

type Token struct {
	Token string `json:"token"`
}

var mux2 *http.ServeMux
var writer *httptest.ResponseRecorder
var users Users
var tokens []Token
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

func TestHandleRegister(t *testing.T) {
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

func TestHandleLogin(t *testing.T) {
	endpoint := "/sso/login"
	u, _ := url.ParseRequestURI(apiURL)
	u.Path = endpoint
	urlPath := u.String()
	for i := 0; i < len(users.Users); i++ {
		jsonData, _ := json.Marshal(map[string]string{
			"email":    users.Users[i].Email,
			"password": users.Users[i].Password,
		})
		request, _ := http.Post(urlPath, "application/json", bytes.NewBuffer(jsonData))
		body, _ := ioutil.ReadAll(request.Body)
		token := Token{}
		json.Unmarshal(body, &token)
		tokens = append(tokens, token)
		if request.StatusCode != 200 {
			t.Errorf("Response code is %v\n", request.StatusCode)
			break
		}
	}
}

func TestHandleVerify(t *testing.T) {
	endpoint := "/sso/verify"
	u, _ := url.ParseRequestURI(apiURL)
	u.Path = endpoint
	urlPath := u.String()
	for _, v := range tokens {
		jsonData, _ := json.Marshal(v)
		request, _ := http.Post(urlPath, "application/json", bytes.NewBuffer(jsonData))
		body, _ := ioutil.ReadAll(request.Body)
		token := Token{}
		json.Unmarshal(body, &token)
		tokens = append(tokens, token)
		if request.StatusCode != 200 {
			t.Errorf("Response code is %v\n", request.StatusCode)
			break
		}
	}
}