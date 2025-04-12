package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Thread struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func GetAllThreadHandler(w http.ResponseWriter, r *http.Request) {
	var threads []Thread

	var query = `SELECT id, title, username, content, created_at FROM threads`
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "failed to query threads", http.StatusInternalServerError)
		log.Printf("failed to query threads: %v", err)
		return
	}

	for rows.Next() {
		var thread Thread
		err := rows.Scan(&thread.ID, &thread.Title, &thread.Username, &thread.Content, &thread.CreatedAt)
		if err != nil {
			http.Error(w, "failed to scan thread row", http.StatusInternalServerError)
			log.Printf("failed to scan thread row: %v", err)
			return
		}
		threads = append(threads, thread)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(threads)
	if err != nil {
		http.Error(w, "Ошибка при кодировании JSON", http.StatusInternalServerError)
	}
}

func GetThreadHandler(w http.ResponseWriter, r *http.Request, id string) {
	var thread Thread

	var query = `SELECT id, title, username, content, created_at FROM threads WHERE id = $1`
	err := db.QueryRow(query, id).Scan(&thread.ID, &thread.Title, &thread.Username, &thread.Content, &thread.CreatedAt)
	if err != nil {
		http.Error(w, "failed to query thread", http.StatusInternalServerError)
		log.Printf("failed to query thread: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(thread)
	if err != nil {
		http.Error(w, "failed encode JSON", http.StatusInternalServerError)
	}
}

func CreateThreadHandler(w http.ResponseWriter, r *http.Request) {
	var thread Thread

	err := json.NewDecoder(r.Body).Decode(&thread)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	thread.CreatedAt = time.Now()

	id, err := CreateThread(&thread)
	if err != nil {
		http.Error(w, "failed to create thread row", http.StatusInternalServerError)
		log.Printf("failed to create thread row: %v", err)
		return
	}

	type Response struct {
		ID int64
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(Response{ID: id})
	if err != nil {
		http.Error(w, "failed encode JSON", http.StatusInternalServerError)
	}
}

func CreateThread(thread *Thread) (int64, error) {
	var query = `INSERT INTO threads (title, username, content, created_at) VALUES ($1, $2, $3, $4)`

	result, err := db.Exec(query, thread.Title, thread.Username, thread.Content, thread.CreatedAt)
	if err != nil {
		return -1, fmt.Errorf("failed to exec query: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("failed to retrieve last insert id: %w", err)
	}

	return id, nil
}
