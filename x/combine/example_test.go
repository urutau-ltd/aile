package combine_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"codeberg.org/urutau-ltd/aile/v2/x/combine"
)

func ExampleMiddleware() {
	var order []string

	h := combine.Middleware(
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, "mw1-before")
				next.ServeHTTP(w, r)
				order = append(order, "mw1-after")
			})
		},
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, "mw2-before")
				next.ServeHTTP(w, r)
				order = append(order, "mw2-after")
			})
		},
	)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "handler")
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	fmt.Println(strings.Join(order, ","))
	// Output: mw1-before,mw2-before,handler,mw2-after,mw1-after
}
