package middlewares

import (
	"net/http"
	"time"

	"forum/database"
)

// CheckSessionCookie validates session cookie and redirects depending on requiresAuth
// true: verifies the user's session cookie before allowing access to protected routes.
// false: prevents already-logged-in users from accessing login/register pages.
func CheckSessionCookie(handler http.HandlerFunc, requiresAuth bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		// 1. cookie existence in request
		if err == nil && cookie.Value != "" {
			var expiryTime time.Time
			err = database.Database.QueryRow(
				"SELECT expires_at FROM sessions WHERE id = ?", cookie.Value,
			).Scan(&expiryTime) // or use "select exists(...)"
			// 2. session existence in database
			if err == nil {
				// 3. valid session expiry datetime
				if expiryTime.After(time.Now()) {
					if requiresAuth {
						handler(w, r)
					} else {
						http.Redirect(w, r, "/", http.StatusSeeOther)
					}
					return
				}
				err = database.DeleteSession(cookie.Value)
				http.SetCookie(w, &http.Cookie{ // all fields needed ?
					Name:     "session_id",
					Value:    "",
					Path:     "/",
					MaxAge:   -1,
					Expires:  time.Now().Add(-1 * time.Hour),
					HttpOnly: true,
				})
			}
		}
		if requiresAuth {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		} else {
			handler(w, r)
		}
	}
}
