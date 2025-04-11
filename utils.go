package main

import "golang.org/x/crypto/bcrypt"

func StartsWith(s, prefix string) bool {
	return (len(s) >= len(prefix)) && (s[:len(prefix)] == prefix)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}
