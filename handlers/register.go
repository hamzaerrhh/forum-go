package handlers

import (
	"net/http"
	"net/mail"
	"regexp"
	"strings"

	"forum/database"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id              int    `json:"id"`   // check for google
	Name            string `json:"name"` // name or username: problem for providers!
	Email           string `json:"email"`
	Password        string
	confirmPassword string
	Message         string
	// Picture string `json:"picture"`    // gmail picture: sometimes cannot be loaded!
	// Avatar  string `json:"avatar_url"` // github avatar
}

func isValidUsername(username string) bool {
	re := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_ ]{1,49}$`)
	// disallowing multiple spaces
	return re.MatchString(username) && !strings.Contains(username, "  ")
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return len(email) >= 5 && len(email) <= 100 && (err == nil)
}

func isValidPassword(password string) bool {
	return len(password) >= 6 && len(password) <= 20
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/register" {
		HandleError(w, http.StatusNotFound, "Page not found")
		return
	}

	switch r.Method {
	case http.MethodGet:
		RenderTemplate(w, 200, "register.html", nil)

	case http.MethodPost:
		user := User{
			Name:            strings.TrimSpace(r.FormValue("name")),
			Email:           strings.TrimSpace(r.FormValue("email")),
			Password:        r.FormValue("password"),
			confirmPassword: r.FormValue("confirm_password"),
		}

		var rules string = `
. username (valid) : 2 ~ 50  chars 
. email (valid)    : 5 ~ 100 chars
. password         : 6 ~ 20  chars`

		// 1. check emptiness
		if user.Name == "" || user.Email == "" || user.Password == "" || user.confirmPassword == "" {
			user.Message = "All fields are required"
			RenderTemplate(w, 400, "register.html", user)
			return
		}
		// 2. check validity
		if !isValidUsername(user.Name) || !isValidEmail(user.Email) || !isValidPassword(user.Password) {
			user.Message = rules
			RenderTemplate(w, 400, "register.html", user)
			return
		}
		// 3. check password match
		if user.Password != user.confirmPassword {
			user.Message = "password and confirm password do not match"
			RenderTemplate(w, 400, "register.html", user)
			return
		}

		// Check email availability
		var emailExists bool
		err := database.Database.QueryRow(
			"SELECT EXISTS(SELECT * FROM users WHERE email = ?)", user.Email,
		).Scan(&emailExists)
		if err != nil {
			HandleError(w, http.StatusInternalServerError, "Database error")
			return
		}
		if emailExists {
			user.Message = "Email already registered" // not good practice
			RenderTemplate(w, 400, "register.html", user)
			return
		}

		// Check username availability
		var nameExists bool
		err = database.Database.QueryRow(
			"SELECT EXISTS(SELECT * FROM users WHERE name = ?)", user.Name,
		).Scan(&nameExists)
		if err != nil {
			HandleError(w, http.StatusInternalServerError, "Database error")
			return
		}
		if nameExists {
			user.Message = "Username already taken"
			RenderTemplate(w, 400, "register.html", user)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			HandleError(w, http.StatusInternalServerError, "Password hashing error")
			return
		}

		_, err = database.Database.Exec(
			"INSERT INTO users (name, email, password) VALUES (?, ?, ?)",
			user.Name,
			user.Email,
			string(hashedPassword),
		)
		// create session if you want to redirect to its page
		if err != nil {
			// log.Println(err.Error())
			HandleError(w, http.StatusInternalServerError, "Could not create account")
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)

	default:
		HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}
