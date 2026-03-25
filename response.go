package aile

import "net/http"

// Status writes a status code with no response body.
func Status(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

// Text writes a plain text HTTP response.
func Text(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(body))
}

// Error writes an HTTP error response.
func Error(w http.ResponseWriter, status int, msg string) {
	http.Error(w, msg, status)
}
