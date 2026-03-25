package basicauth

import (
	"net/http"

	"codeberg.org/urutau-ltd/aile/v2"
)

// Validator checks a username and password pair.
type Validator func(user, pass string) bool

// Middleware enforces HTTP Basic authentication for the wrapped handler.
func Middleware(realm string, validate Validator) aile.Middleware {
	if realm == "" {
		realm = "restricted"
	}
	challenge := `Basic realm="` + realm + `"`

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, pass, ok := r.BasicAuth()

			if !ok || validate == nil || !validate(user, pass) {
				w.Header().Set("WWW-Authenticate", challenge)
				aile.Error(w,
					http.StatusUnauthorized,
					http.StatusText(http.StatusUnauthorized))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
