package handlers

import (
	"log"
	"net/http"

	"forum/database"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/logout" {
		HandleError(w, http.StatusNotFound, "Page not found")
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil { // http.ErrNoCookie
		return
	}

	err = deleteSession(cookie.Value)
	// + need to remove cookie from storage
	if err != nil {
		log.Println(err)
	}
	http.SetCookie(w, &http.Cookie{Name: "session_id", MaxAge: -1}) // delete cookie
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func deleteSession(sessionId string) error {
	// return database.Database.QueryRow(query).Err()
	// db.exec vs db.queryrow in golang sqlite
	// queryrow not working with delete statement
	_, err := database.Database.Exec(
		"DELETE FROM sessions WHERE id = ?",
		sessionId) // returns result
	return err
}
