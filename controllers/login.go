package controllers

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type LoginData struct {
	Err string
}

var users = map[string]string{
	"admin": "admin",
}

var tokens = map[string]string{}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	loginpage(&LoginData{}, w, r)
}

func LoginAction(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form["username"]
	password := r.Form["password"]

	log.Printf("Login request with %s and %s\n", username, password)

	if len(username) == 0 || username[0] == "" || users[username[0]] == "" {
		loginpage(&LoginData{
			Err: "invalid login request",
		}, w, r)
		log.Println("invalid username")
		return
	}

	if len(password) == 0 || password[0] == "" || users[username[0]] != password[0] {
		loginpage(&LoginData{
			Err: "invalid login request",
		}, w, r)
		log.Println("invalid password")
		return
	}

	log.Println("login OK")
	token := uuid.New().String()
	tokens[token] = username[0]
	http.SetCookie(w, &http.Cookie{
		Name:     "id",
		Value:    token,
		HttpOnly: true,
		Expires:  time.Now().Add(time.Hour * 1),
	})

	http.Redirect(w, r, r.URL.Scheme+r.URL.Host+"/home", 302)
	w.Write([]byte("Logged In welocme"))
}

func loginpage(data *LoginData, w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("tmpls/login.html")
	if err != nil {
		log.Println(err)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Println(err)
	}
}

func GetUsername(token string) string {
	return tokens[token]
}

func CheckUserAndRespond(w http.ResponseWriter, r *http.Request) bool {
	se_cookie, err := r.Cookie("id")
	if err != nil || se_cookie.Value == "" || GetUsername(se_cookie.Value) == "" {
		w.WriteHeader(http.StatusForbidden)
		return false
	}
	return true
}
