package main

import (
	"fmt"
	"net/http"

	"zone/database"
	zone "zone/handlers"
)

func main() {
	database.Init()
	http.HandleFunc("/", zone.Home)
	http.HandleFunc("/Register", zone.Register)
	http.HandleFunc("/Login", zone.Login)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
