package compress_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	aile "codeberg.org/urutau-ltd/aile/v2"
	"codeberg.org/urutau-ltd/aile/v2/x/compress"
)

func ExampleMiddleware() {
	h := compress.Middleware(compress.Config{MinSize: 1})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		aile.Text(w, http.StatusOK, "hello compressed world")
	}))

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	fmt.Println(rec.Header().Get("Content-Encoding"))
	// Output: gzip
}
