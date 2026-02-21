package main

import (
	"fmt"
	"net/http"

	"forum/database"
	"forum/handlers"
)

func main() {
	if err := database.Init(); err != nil {
		fmt.Println(err)
	}

	http.HandleFunc("/", handlers.Forum)
	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/logout", handlers.Logout)

	http.HandleFunc("/static/styles.css", handlers.Styles)

	fmt.Println("Server running on http://0.0.0.0:8080")
	http.ListenAndServe(":8080", nil)
}
