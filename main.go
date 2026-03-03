package main

import (
	"fmt"
	"forum/database"
	"forum/handlers"
	"forum/routing"
	"log"
	"net/http"
)

func main() {
	if err := database.Init(); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	http.HandleFunc("/static/", handlers.HandleStatic)
	http.HandleFunc("/", handlers.Forum) // use middleware when separated to home & feed

	routing.RegisterRoutes()
	// Static
	// http.HandleFunc("/static/", zone.HandleStatic)

	fmt.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

