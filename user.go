package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID       int64
	Username string
	Password string
}

func GenerateSessionToken() (string, error) {
	const length = 64
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func CreateUser(user *User) error {
	var query = `insert into users (username, password) values ($1, $2)`

	_, err := db.Exec(query, user.Username, user.Password)
	if err != nil {
		return fmt.Errorf("failed to exec query: %w", err)
	}
	return nil
}

func GetUserByUsername(username string, user *User) error {
	var query = "select id, username, password FROM users WHERE username = $1"
	return db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password)
}

func UserSigninHandler(w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("Username")

	var user User
	if err := GetUserByUsername(username, &user); err != nil {
		log.Fatal("failed to get user: %w", err)
		return
	}

	/* TODO(vlad0924): change to password_hash and error message */
	password := r.PostFormValue("Password")
	if user.Password != password {
		log.Fatal("password invalid")
		return
	}

	token, err := GenerateSessionToken()
	if err != nil {
		log.Fatal("failed get generate token: %w", err)
		return
	}
	expiry := time.Now().Add(OneWeek)

	session := &Session{
		ID:     user.ID,
		Expiry: expiry,
		User:   user,
	}

	Sessions[token] = session

	SetCookieToken(w, r, token)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func UserSignupHandler(w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("Username")
	password := r.PostFormValue("Password")

	// password_hash, err := HashPassword(password)
	// if err != nil {
	// 	log.Fatal("failed transform password to hash: %v", err)
	// 	return
	// }

	user := User{
		Username: username,
		Password: password,
	}
	CreateUser(&user)
	http.Redirect(w, r, "/signin", http.StatusSeeOther)
}

func UserSigninPageHandler(w http.ResponseWriter, r *http.Request) {
	// 	const pageFormat = `
	// <!DOCTYPE html>
	// <head>
	// 	<title>Sign in</title>
	// </head>
	// <body>
	// 	<h1>Web Forum</h1>
	// 	<h2>Sign in</h2>
	// 	<form method="POST" action="/api/user/signin">
	// 		<label>Username:
	// 			<input type="username" name="Username" required>
	// 		</label>
	// 		<br></br>
	// 		<label>Password:
	// 			<input type="password" name="Password" required>
	// 		</label>
	// 		<br></br>
	// 		<input type="submit" value="Sign in">
	// 	</form>
	// </body>
	// </html>
	// `

	// 	w.Header().Add("Content-Type", "text/html")
	// 	w.Write([]byte(pageFormat))

	t, _ := template.ParseFiles("html/signin.html")
	if err := t.Execute(w, nil); err != nil {
		log.Println("failed to load signin page")
	}
}

func UserSignupPageHandler(w http.ResponseWriter, r *http.Request) {
	// 	const pageFormat = `
	// <!DOCTYPE html>
	// <head>
	// 	<title>Sign in</title>
	// </head>
	// <body>
	// 	<h1>Web Forum</h1>
	// 	<h2>Sign up</h2>
	// 	<form method="POST" action="/api/user/signup">
	// 		<label>Username:
	// 			<input type="username" name="Username" required>
	// 		</label>
	// 		<br></br>
	// 		<label>Password:
	// 			<input type="password" name="Password" required>
	// 		</label>
	// 		<br></br>
	// 		<input type="submit" value="Sign up">
	// 	</form>
	// </body>
	// </html>
	// `

	// 	w.Header().Add("Content-Type", "text/html")
	// 	w.Write([]byte(pageFormat))

	t, _ := template.ParseFiles("html/signup.html")
	if err := t.Execute(w, nil); err != nil {
		log.Println("failed to load signup page")
	}
}
