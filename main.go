package main

import (
	"log"
	"net/http"
)

func AuthMiddleware(r *http.Request) bool {
	cookie, err := r.Cookie("Token")
	if err != nil {
		log.Println("failed to get cookie")
	}

	_, err = GetSessionFromToken(cookie.Value)
	return err == nil
}

func HandleAPIRequest(w http.ResponseWriter, r *http.Request, path string) {
	switch {
	case StartsWith(path, "/user"):
		switch path[len("/user"):] {
		case "/signin":
			UserSigninHandler(w, r)
		case "/signup":
			UserSignupHandler(w, r)
		}
	case StartsWith(path, "/thread"):
		if AuthMiddleware(r) {
			if path[len("/thread"):] == "" {
				if r.Method == "GET" {
					GetAllThreadHandler(w, r)
				} else if r.Method == "POST" {
					CreateThreadHandler(w, r)
				}
			} else if (len(path) > len("/thread/")) && (path[:len("/thread/")] == "/thread/") { // /thread/{threadID}
				if r.Method == "GET" {
					GetThreadHandler(w, r, path[len("/thread/"):])
				}
			} else {
				// return error HTTPStatusNotFound
			}
		} else {
			log.Println("access denied")
		}
	}
}

type Router struct{}

func (rt *Router) RouterFunc(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	switch {
	case StartsWith(path, "/api"):
		HandleAPIRequest(w, r, path[len("/api"):])
	}
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rt.RouterFunc(w, r)
}

func main() {
	if err := OpenDB("db"); err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer CloseDB()

	router := &Router{}
	log.Println("Server start...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
