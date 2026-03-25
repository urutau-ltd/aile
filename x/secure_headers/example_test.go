package secureheaders_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	aile "codeberg.org/urutau-ltd/aile/v2"
	secureheaders "codeberg.org/urutau-ltd/aile/v2/x/secure_headers"
)

func ExampleMiddleware() {
	h := secureheaders.Middleware(secureheaders.Config{
		ContentTypeNosniff: true,
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		aile.Status(w, http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	fmt.Println(rec.Header().Get("X-Content-Type-Options"))
	// Output: nosniff
}
