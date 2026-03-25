package bearerauth_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	aile "codeberg.org/urutau-ltd/aile/v2"
	bearerauth "codeberg.org/urutau-ltd/aile/v2/x/bearer_auth"
)

func ExampleMiddleware() {
	h := bearerauth.Middleware(func(token string) bool {
		return token == "good-token"
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		aile.Status(w, http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/private", nil)
	req.Header.Set("Authorization", "Bearer good-token")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	fmt.Println(rec.Code)
	// Output: 204
}
