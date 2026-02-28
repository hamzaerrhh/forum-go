package handlers

import (
	"bytes"
	"html/template"
	"net/http"
)

func RenderTemplate(w http.ResponseWriter, status int, tmpl string, data any) {
	t, err := template.ParseFiles("templates/" + tmpl)
	if err != nil {
		HandleError(w, http.StatusInternalServerError, "Template error")
		return
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		HandleError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}