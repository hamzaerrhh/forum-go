package main

import (
	"fmt"
	"log"
	"net/http"

	"forum/database"
	"forum/handlers"
	"forum/middlewares"
)

func main() {
	if err := database.Init(); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	http.HandleFunc("/", handlers.Forum) // use middleware when separated to home & feed
	http.HandleFunc("/register", middlewares.CheckSessionCookie(handlers.Register, false))
	http.HandleFunc("/login", middlewares.CheckSessionCookie(handlers.Login, false))
	http.HandleFunc("/logout", middlewares.CheckSessionCookie(handlers.Logout, true))

	http.HandleFunc("/static/styles.css", handlers.Styles)
	// http.HandleFunc("/static/", zone.HandleStatic)

	fmt.Println("Server running on http://0.0.0.0:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
