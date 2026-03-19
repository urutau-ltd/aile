package aile

import (
	"encoding/json"
	"net/http"
)

// Writes a JSON HTTP Response. It already sets the Content-Type
// header and the charset is set to uft-8. Mind this to avoid
// header duplication in your code.
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// Decode JSON request body into a given type T.
func DecodeJSON[T any](r *http.Request) (T, error) {
	var v T
	err := json.NewDecoder(r.Body).Decode(&v)
	return v, err
}
