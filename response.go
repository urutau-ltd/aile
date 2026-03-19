package aile

import "net/http"

// Writes a status code with no request body.
func Status(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

// Writes a Plain Text HTTP Response
func Text(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(body))
}

// Writes an HTTP error response.
func Error(w http.ResponseWriter, status int, msg string) {
	http.Error(w, msg, status)
}
