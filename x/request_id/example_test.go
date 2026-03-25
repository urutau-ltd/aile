package requestid_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	aile "codeberg.org/urutau-ltd/aile/v2"
	requestid "codeberg.org/urutau-ltd/aile/v2/x/request_id"
)

func ExampleMiddleware() {
	h := requestid.Middleware(requestid.Config{
		Header: "X-Request-ID",
		Generator: func() string {
			return "req-1"
		},
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, _ := requestid.FromContext(r.Context())
		fmt.Println(id)
		aile.Status(w, http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	fmt.Println(rec.Header().Get("X-Request-ID"))
	// Output:
	// req-1
	// req-1
}
