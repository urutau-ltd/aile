package aile

import (
	"encoding/json"
	"net/http"
)

// WriteJSON writes a JSON response and sets the Content-Type header to
// application/json with utf-8 charset.
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// DecodeJSON decodes the JSON request body into a value of type T.
func DecodeJSON[T any](r *http.Request) (T, error) {
	var v T
	err := json.NewDecoder(r.Body).Decode(&v)
	return v, err
}
