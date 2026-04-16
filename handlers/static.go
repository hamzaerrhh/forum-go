package handlers

import (
	"net/http"
	"os"
	"path/filepath"
)

// HandleStatic serves CSS files and prevents direct access to /static/
// func Static(w http.ResponseWriter, r *http.Request) {
// 	if r.URL.Path == "/static" || r.URL.Path == "/static/" {
// 		HandleError(w, http.StatusNotFound, "Not Found!")
// 		return
// 	}
// 	filePath := filepath.Join("assets", r.URL.Path[len("/static/"):])
// 	if _, err := os.Stat(filePath); err != nil {
// 		HandleError(w, http.StatusNotFound, "Not Found!")
// 		return
// 	}
// 	http.ServeFile(w, r, filePath)
// }

func Static(w http.ResponseWriter, r *http.Request) {
	var baseDir string

	switch {
	case r.URL.Path[:8] == "/static/":
		baseDir = "static"
		// remove "/static/" prefix
		r.URL.Path = r.URL.Path[8:]
	case r.URL.Path[:9] == "/uploads/":
		baseDir = "uploads"
		// remove "/uploads/" prefix
		r.URL.Path = r.URL.Path[9:]
	default:
		http.NotFound(w, r)
		return
	}

	if r.URL.Path == "" {
		http.NotFound(w, r)
		return
	}

	filePath := filepath.Join(baseDir, r.URL.Path)
	if _, err := os.Stat(filePath); err != nil {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, filePath)
}
