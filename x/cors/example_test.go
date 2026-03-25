package cors_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	aile "codeberg.org/urutau-ltd/aile/v2"
	"codeberg.org/urutau-ltd/aile/v2/x/cors"
)

func ExampleMiddleware() {
	h := cors.Middleware(cors.Config{
		AllowOrigins: []string{"https://app.example.com"},
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		aile.Status(w, http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Origin", "https://app.example.com")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	fmt.Println(rec.Header().Get("Access-Control-Allow-Origin"))
	// Output: https://app.example.com
}
