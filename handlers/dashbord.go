package zone

import (
	"html/template"
	"net/http"

	"zone/database"
)


func Dashboard(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}

	var user User
	err = database.Database.QueryRow(
		"SELECT userName, email, password FROM users WHERE session = ? AND dateexpired > DATETIME('now')",
		cookie.Value,
	).Scan(&user.Name, &user.Email, &user.Password)
	if err != nil {
		http.Redirect(w, r, "/Login", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("templates/Dashboard.html")
	if err != nil {
		HandleError(w, 500, "Template error")
		return
	}

	tmpl.Execute(w, user)
}
