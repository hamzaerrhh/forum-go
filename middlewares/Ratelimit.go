package middlewares

import (
	"database/sql"
	"net"
	"net/http"
	"time"

	"forum/database"
	"forum/handlers"
)

func RateLimit(handler http.HandlerFunc, maxRequests int, window time.Duration) http.HandlerFunc {
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

		var count int // number of requests
		var lastRequest time.Time

		err = database.Database.QueryRow(
			"SELECT count, last_request FROM rate_limits WHERE ip = ?", ip,
		).Scan(&count, &lastRequest)

		if err == sql.ErrNoRows {
			_, err = database.Database.Exec(
				"INSERT INTO rate_limits (ip, count, last_request) VALUES (?, ?, ?)",
				ip, 1, time.Now(),
			)
			if err != nil {
				handlers.HandleError(w, http.StatusInternalServerError, "Internal Server Erro")
				return
			}
			handler(w, r)
			return
		} else if err != nil {
			handlers.HandleError(w, http.StatusInternalServerError, "Internal Server Erro")
			return
		}

		if time.Since(lastRequest) > window {
			count = 0
		}

		if count >= maxRequests {
			handlers.HandleError(w, http.StatusTooManyRequests, "Rate limit exceeded. Please try again later.")
			return
		}

		_, err = database.Database.Exec(
			"UPDATE rate_limits SET count = ?, last_request = ? WHERE ip = ?",
			count+1, time.Now(), ip,
		)
		if err != nil {
			handlers.HandleError(w, http.StatusInternalServerError, "Internal Server Erro")
			return
		}

		handler(w, r)
	}
}
