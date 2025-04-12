package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func GenerateSessionToken() (string, error) {
	const length = 10
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func UserDecodeJSONInput(w http.ResponseWriter, r *http.Request, data interface{}) {
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func CreateUser(user *User) (int64, error) {
	var err error
	user.Password = HashPassword(user.Password)

	var query = `INSERT INTO users (username, password) values ($1, $2)`

	result, err := db.Exec(query, user.Username, user.Password)
	if err != nil {
		return -1, fmt.Errorf("failed to exec query: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("failed to retrieve last insert id: %w", err)
	}

	return id, nil
}

func GetUserByUsername(username string, user *User) error {
	var query = "SELECT id, username, password FROM users WHERE username = $1"
	return db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password)
}

type userSigninInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func UserSigninHandler(w http.ResponseWriter, r *http.Request) {
	var input userSigninInput

	UserDecodeJSONInput(w, r, &input)

	var user User
	if err := GetUserByUsername(input.Username, &user); err != nil {
		log.Fatal("failed to get user: %w", err)
		return
	}

	if HashPassword(input.Password) != user.Password {
		log.Fatal("failed qi passwords")
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

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(tokenResponse{Token: token})
	if err != nil {
		http.Error(w, "failed encode JSON", http.StatusInternalServerError)
	}
}

type userSignupInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func UserSignupHandler(w http.ResponseWriter, r *http.Request) {
	var input userSignupInput

	UserDecodeJSONInput(w, r, &input)

	user := User{
		Name:     input.Name,
		Email:    input.Email,
		Username: input.Username,
		Password: input.Password,
	}

	id, err := CreateUser(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(idResponse{ID: id})
	if err != nil {
		http.Error(w, "failed encode JSON", http.StatusInternalServerError)
	}
}
