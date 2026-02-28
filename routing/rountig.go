package routing

import (
	"net/http"
	"time"

	"forum/handlers"
	"forum/middlewares"
)

func RegisterRoutes() {

	http.HandleFunc(
		"/posts/create",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.CreatePost, true),
			5,
			time.Minute,
		),
	)

	http.HandleFunc(
		"/comments/create",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.CreateComment, true),
			10,
			time.Minute,
		),
	)

	http.HandleFunc(
		"/login",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.Login, false),
			3,
			time.Minute,
		),
	)

	http.HandleFunc(
		"/register",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.Register, false),
			3,
			time.Minute,
		),
	)

	http.HandleFunc(
		"/logout",
		middlewares.RateLimit(
			middlewares.CheckSessionCookie(handlers.Logout, true),
			5,
			time.Minute,
		),
	)
}