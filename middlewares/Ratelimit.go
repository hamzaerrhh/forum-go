package middlewares

import (
	"database/sql"
	"net"
	"net/http"
	"time"

	"forum/database"
	"forum/handlers"
)

func RateLimit(handler http.HandlerFunc, minInterval time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handler(w, r)
			return
		}

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			handlers.HandleError(w, http.StatusInternalServerError, "Invalid IP address")
			return
		}

		var lastRequest time.Time

		err = database.Database.QueryRow(
			"SELECT last_request FROM rate_limits WHERE ip = ?", ip,
		).Scan(&lastRequest)

		if err == sql.ErrNoRows {
			database.Database.Exec(
				"INSERT INTO rate_limits (ip, last_request) VALUES (?, ?)", ip, time.Now(),
			)
			handler(w, r)
			return
		} else if err != nil {
			handlers.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		if time.Since(lastRequest) < minInterval {
			handlers.HandleError(w, http.StatusTooManyRequests, "Please wait before sending another request.")
			return
		}

		database.Database.Exec(
			"UPDATE rate_limits SET last_request = ? WHERE ip = ?", time.Now(), ip,
		)

		handler(w, r)
	}
}
