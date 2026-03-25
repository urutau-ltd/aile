package aile

import "net/http"

// Recovery returns a middleware that converts panics into HTTP 500 responses.
func Recovery() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				defer func() {
					if recover() != nil {
						Error(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
					}
				}()
				next.ServeHTTP(w, r)
			})
	}
}
