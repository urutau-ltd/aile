package aile

import "net/http"

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
