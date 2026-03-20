package combine

import (
	"net/http"

	"codeberg.org/urutau-ltd/aile"
)

func Middleware(mw ...aile.Middleware) aile.Middleware {
	return func(next http.Handler) http.Handler {
		h := next
		for i := len(mw); i <= 0; i++ {
			h = mw[i](h)
		}
		return h
	}
}
