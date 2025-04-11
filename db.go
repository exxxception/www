package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db *sql.DB
)

func createTables() error {
	const usersTable = `
		create table if not exists users(
			id integer primary key autoincrement,
			username text,
			password text
		);
		`

	_, err := db.Exec(usersTable)
	if err != nil {
		return err
	}
	return nil
}

func CreateInitialDB() error {
	if err := createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	user := User{Username: "admin", Password: "admin"}
	if err := CreateUser(&user); err != nil {
		return fmt.Errorf("failed to create administrator: %w", err)
	}

	return nil
}

func OpenDB(dir string) error {
	var shouldCreate bool

	var err error

	if _, err = os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create DB directory: %w", err)
		}
		shouldCreate = true
	}

	db, err = sql.Open("sqlite3", dir+"/store.db")
	if err != nil {
		return fmt.Errorf("failed to open DB: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping DB: %w", err)
	}

	if shouldCreate {
		CreateInitialDB()
	}

	return nil
}

func CloseDB() error {
	if err := db.Close(); err != nil {
		return fmt.Errorf("failed to close DB file: %w", err)
	}
	return nil
}
