package main

import (
	"errors"
	"net/http"
	"time"
)

type Session struct {
	ID     int64
	Expiry time.Time
	User   User
}

const OneWeek = time.Hour * 24 * 7

var (
	Sessions = make(map[string]*Session)
)

func GetSessionFromToken(token string) (*Session, error) {
	session, ok := Sessions[token]
	if !ok {
		return nil, errors.New("session for this token does not exist")
	}
	return session, nil
}

func SetCookieToken(w http.ResponseWriter, r *http.Request, token string) {
	cookie := http.Cookie{
		Name:     "Token",
		Value:    token,
		Path:     "/",
		MaxAge:   int(time.Now().Add(OneWeek).Unix()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)
}
