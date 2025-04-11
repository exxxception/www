package main

import (
	"fmt"
	"log"
	"net/http"
)

func HandlePageRequest(w http.ResponseWriter, r *http.Request, path string) {
	switch path {
	case "/":
		IndexPageHandler(w, r)
	case "/signin":
		UserSigninPageHandler(w, r)
	case "/signup":
		UserSignupPageHandler(w, r)
	}
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
	}
}

type Router struct{}

func (rt *Router) RouterFunc(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	switch {
	default:
		HandlePageRequest(w, r, path)
	case StartsWith(path, "/api"):
		HandleAPIRequest(w, r, path[len("/api"):])
	}
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rt.RouterFunc(w, r)
}

func IndexPageHandler(w http.ResponseWriter, r *http.Request) {
	var indexPage = `
<!DOCTYPE html>
<head>
	<title>Web Forum</title>
</head>
<body>
	<h1>Welcome</h1>
	<a href="/signin">Sign in</a>
	<a href="/signup">Sign up</a>
</body>
</html>
`

	row, err := r.Cookie("Token")
	if err == nil {
		session, err := GetSessionFromToken(row.Value)
		if err == nil {
			indexPage = fmt.Sprintf(`
<!DOCTYPE html>
<head>
	<title>Web Forum</title>
</head>
<body>
	<h1>Welcome, %s</h1>
	<a href="/">Logout</a>
</body>
</html>
`, session.User.Username)
		}
	}

	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte(indexPage))
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
