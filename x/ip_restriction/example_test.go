package iprestriction_test

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"

	aile "codeberg.org/urutau-ltd/aile/v2"
	iprestriction "codeberg.org/urutau-ltd/aile/v2/x/ip_restriction"
)

func ExampleMiddleware() {
	_, allow, err := net.ParseCIDR("127.0.0.0/8")
	if err != nil {
		panic(err)
	}

	h := iprestriction.Middleware(iprestriction.Config{
		Allow: []*net.IPNet{allow},
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		aile.Status(w, http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/internal", nil)
	req.RemoteAddr = "127.0.0.1:1234"
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	fmt.Println(rec.Code)
	// Output: 204
}
