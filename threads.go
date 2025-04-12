package main

import (
	"log"
	"net/http"
	"text/template"
	"time"
)

type Thread struct {
	ID       int64
	Author   string
	Avatar   string
	Username string
	Date     string
	Text     string
}

func GetThread(w http.ResponseWriter, r *http.Request) {
	threads := []Thread{
		{
			Author:   "Vlad Korovkin",
			Avatar:   "/html/photo.jpg",
			Username: "vlad",
			Date:     time.Now().Format(time.RFC822),
			Text:     "Всем привет! Помогите закопать яму. Плачу 1000 руб.",
		},
		{
			Author:   "Nikita Lomtik",
			Avatar:   "",
			Username: "bad_boy44",
			Date:     time.Now().Format(time.RFC822),
			Text:     "А вот и второй тред.",
		},
		{
			Author:   "VOlodya",
			Avatar:   "",
			Username: "hacker228",
			Date:     time.Now().Format(time.RFC822),
			Text:     "I'am pro programmer GOlang",
		},
	}
	t, _ := template.ParseFiles("html/index.html", "html/thread.html")
	if err := t.Execute(w, struct{ Threads []Thread }{Threads: threads}); err != nil {
		log.Println("failed to load thread")
	}
}
