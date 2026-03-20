package aile

import (
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"testing"
)

func debugServeMuxRequest(t *testing.T, mux *http.ServeMux, method, path string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(method, path, nil)

	h, pattern := mux.Handler(req)

	t.Logf("go version: %s", runtime.Version())
	t.Logf("GOOS=%s GOARCH=%s", runtime.GOOS, runtime.GOARCH)
	t.Logf("GODEBUG=%q", os.Getenv("GODEBUG"))
	t.Logf("request method=%q path=%q", req.Method, req.URL.Path)
	t.Logf("resolved pattern=%q handler=%T", pattern, h)

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	t.Logf("response status=%d body=%q", rec.Code, rec.Body.String())
	return rec
}
