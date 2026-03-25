package requestid

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"codeberg.org/urutau-ltd/aile/v2"
)

func TestMiddlewarePropagatesIncomingRequestID(t *testing.T) {
	h := Middleware(Config{Header: "X-Request-ID"})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := FromContext(r.Context())
		if !ok || id != "req-123" {
			t.Fatalf("unexpected request id from context: %q %v", id, ok)
		}
		aile.Status(w, http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("X-Request-ID", "req-123")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if got := rec.Header().Get("X-Request-ID"); got != "req-123" {
		t.Fatalf("unexpected response request id: got %q", got)
	}
}

func TestMiddlewareGeneratesRequestIDWhenMissing(t *testing.T) {
	h := Middleware(Config{
		Header: "X-Request-ID",
		Generator: func() string {
			return "generated-id"
		},
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := FromContext(r.Context())
		if !ok || id != "generated-id" {
			t.Fatalf("unexpected request id from context: %q %v", id, ok)
		}
		aile.Status(w, http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if got := rec.Header().Get("X-Request-ID"); got != "generated-id" {
		t.Fatalf("unexpected response request id: got %q", got)
	}
}
