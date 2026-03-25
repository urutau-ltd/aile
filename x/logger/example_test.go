package logger_test

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"

	aile "codeberg.org/urutau-ltd/aile/v2"
	xlogger "codeberg.org/urutau-ltd/aile/v2/x/logger"
)

func ExampleMiddleware() {
	var buf bytes.Buffer

	h := xlogger.Middleware(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, attr slog.Attr) slog.Attr {
			if attr.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return attr
		},
	})))(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		aile.Status(w, http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	fmt.Println(strings.Contains(buf.String(), "status=204"))
	// Output: true
}
