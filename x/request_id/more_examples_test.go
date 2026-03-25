package requestid_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	requestid "codeberg.org/urutau-ltd/aile/v2/x/request_id"
)

func ExampleConfig() {
	cfg := requestid.Config{
		Header: "X-Request-ID",
	}

	fmt.Println(cfg.Header)
	// Output: X-Request-ID
}

func ExampleFromContext() {
	h := requestid.Middleware(requestid.Config{
		Generator: func() string { return "req-123" },
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := requestid.FromContext(r.Context())
		fmt.Println(id, ok)
	}))

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	// Output: req-123 true
}
