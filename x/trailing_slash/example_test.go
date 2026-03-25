package trailingslash_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	trailingslash "codeberg.org/urutau-ltd/aile/v2/x/trailing_slash"
)

func ExampleMiddleware() {
	h := trailingslash.Middleware(trailingslash.RedirectTrim)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/users/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	fmt.Println(rec.Header().Get("Location"))
	// Output: /users
}
