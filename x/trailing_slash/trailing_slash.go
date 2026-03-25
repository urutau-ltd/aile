package trailingslash

import (
	"net/http"
	"strings"

	"codeberg.org/urutau-ltd/aile/v2"
)

// Mode controls how the middleware normalizes trailing slashes.
type Mode int

const (
	// RedirectTrim redirects "/path/" to "/path".
	RedirectTrim Mode = iota + 1
	// RedirectAppend redirects "/path" to "/path/".
	RedirectAppend
)

// Middleware redirects requests according to the selected trailing slash mode.
func Middleware(mode Mode) aile.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			switch mode {
			case RedirectTrim:
				if path != "/" && strings.HasSuffix(path, "/") {
					url := *r.URL
					url.Path = strings.TrimRight(path, "/")
					if url.Path == "" {
						url.Path = "/"
					}
					http.Redirect(w, r, url.String(), http.StatusPermanentRedirect)
					return
				}
			case RedirectAppend:
				if path != "/" && !strings.HasSuffix(path, "/") {
					url := *r.URL
					url.Path = path + "/"
					http.Redirect(w, r, url.String(), http.StatusPermanentRedirect)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
