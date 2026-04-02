package handlers

import (
	"fmt"
	"net/http"
	"net/url"
)

var (
	// Google
	GOOGLE_CLIENT_ID     string
	GOOGLE_CLIENT_SECRET string
	redirectURI          = "http://localhost:8080/auth/google/callback"

	// Github
	GITHUB_CLIENT_ID     string
	GITHUB_CLIENT_SECRET string
)

func OAuthLoginHandler(w http.ResponseWriter, r *http.Request) {
	var baseURL, client_id, scope string

	provider := r.PathValue("provider")
	switch provider {
	case "google":
		baseURL = "https://accounts.google.com/o/oauth2/v2/auth"
		client_id = GOOGLE_CLIENT_ID
		scope = "openid email profile"

	case "github":
		baseURL = "https://github.com/login/oauth/authorize"
		client_id = GITHUB_CLIENT_ID
		scope = "read:user user:email"

	default:
		HandleError(w, http.StatusNotFound, "Unknown endpoint")
	}

	// Generate a random state token to prevent CSRF
	// state := generateState()
	// stateStore[state] = true

	params := url.Values{}
	params.Add("client_id", client_id)

	// only for google ?!
	if provider == "google" {
		params.Add("redirect_uri", redirectURI)
	}

	params.Add("response_type", "code")
	params.Add("scope", scope)

	// only for google ?!
	params.Add("access_type", "offline")

	// params.Add("prompt", "consent")
	// Forces account chooser to appear every time (optional)
	params.Set("prompt", "select_account")

	// params.Set("state", state)

	authURL := baseURL + "?" + params.Encode()

	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func OAuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	var tokenURL, client_id, client_secret, userInfoURL string

	provider := r.PathValue("provider")
	switch provider {
	case "google":
		tokenURL = "https://oauth2.googleapis.com/token"
		client_id = GOOGLE_CLIENT_ID
		client_secret = GOOGLE_CLIENT_SECRET
		userInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"

	case "github":
		tokenURL = "https://github.com/login/oauth/access_token"
		client_id = GITHUB_CLIENT_ID
		client_secret = GITHUB_CLIENT_SECRET
		userInfoURL = "https://api.github.com/user"

	default:
		HandleError(w, http.StatusNotFound, "Unknown endpoint")
	}

	// 1. Validate state to prevent CSRF
	// state := r.URL.Query().Get("state")
	// if !stateStore[state] {
	// 	http.Error(w, "invalid state", http.StatusBadRequest)
	// 	return
	// }
	// delete(stateStore, state) // one-time use

	// 2. Check for errors (e.g. user denied consent)
	if errMsg := r.URL.Query().Get("error"); errMsg != "" {
		http.Error(w, "auth error: "+errMsg, http.StatusUnauthorized)
		return
	}

	// 3. Exchange authorization code for tokens
	code := r.URL.Query().Get("code")
	tokenData, err := exchangeCode(provider, tokenURL, client_id, client_secret, code)
	if err != nil {
		http.Error(w, "token exchange failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	accessToken := tokenData.AccessToken

	// 4. Fetch user info
	user, err := fetchUserInfo(userInfoURL, accessToken)
	if provider == "github" {
		user.Email, _ = fetchUserEmail(accessToken)
	}
	if err != nil {
		http.Error(w, "failed to fetch user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 5. At this point, create your session / JWT / cookie
	// For demo, just display the user info
	w.Header().Set("Content-Type", "text/html")
	if provider == "google" {
		fmt.Fprintf(w, `<h2>Logged in!</h2>
			<p>Name: %s</p>
			<p>Email: %s</p>
			<img src="%s">`,
			user.Name, user.Email, user.Picture)
	} else {
		fmt.Fprintf(w, `<h2>Logged in!</h2>
			<p>Name: %s</p>
			<p>Email: %s</p>
			<img src="%s">`,
			user.Name, user.Email, user.Avatar)
	}
}
