package basicauth_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	aile "codeberg.org/urutau-ltd/aile/v2"
	basicauth "codeberg.org/urutau-ltd/aile/v2/x/basic_auth"
)

func ExampleMiddleware() {
	h := basicauth.Middleware("admin", func(user, pass string) bool {
		return user == "admin" && pass == "secret"
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		aile.Status(w, http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.SetBasicAuth("admin", "secret")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	fmt.Println(rec.Code)
	// Output: 204
}
