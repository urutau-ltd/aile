package bearerauth

import (
	"net/http"
	"strings"

	"codeberg.org/urutau-ltd/aile/v2"
)

// Validator checks whether a bearer token is accepted.
type Validator func(token string) bool

// Middleware enforces Bearer token authentication for the wrapped handler.
func Middleware(validate Validator) aile.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authz := r.Header.Get("Authorization")
			const prefix string = "Bearer "
			if !strings.HasPrefix(authz, prefix) || validate == nil || !validate(strings.TrimPrefix(authz, prefix)) {
				w.Header().Set("WWW-Authenticate", "Bearer")
				aile.Error(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
