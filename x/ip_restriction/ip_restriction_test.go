package iprestriction

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"codeberg.org/urutau-ltd/aile/v2"
)

func mustCIDR(t *testing.T, s string) *net.IPNet {
	t.Helper()

	_, network, err := net.ParseCIDR(s)
	if err != nil {
		t.Fatalf("ParseCIDR(%q) returned error: %v", s, err)
	}
	return network
}

func TestMiddlewareAllowsMatchingClientIP(t *testing.T) {
	h := Middleware(Config{
		Allow: []*net.IPNet{mustCIDR(t, "127.0.0.0/8")},
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		aile.Status(w, http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.RemoteAddr = "127.0.0.1:1234"
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("unexpected status: got %d want %d", rec.Code, http.StatusNoContent)
	}
}

func TestMiddlewareUsesForwardedIPWhenTrusted(t *testing.T) {
	h := Middleware(Config{
		Allow:      []*net.IPNet{mustCIDR(t, "10.0.0.0/8")},
		TrustProxy: true,
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		aile.Status(w, http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.RemoteAddr = "127.0.0.1:1234"
	req.Header.Set("X-Forwarded-For", "10.1.2.3, 127.0.0.1")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("unexpected status: got %d want %d", rec.Code, http.StatusNoContent)
	}
}

func TestMiddlewareRejectsDeniedClientIP(t *testing.T) {
	h := Middleware(Config{
		Deny: []*net.IPNet{mustCIDR(t, "127.0.0.0/8")},
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		aile.Status(w, http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.RemoteAddr = "127.0.0.1:1234"
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("unexpected status: got %d want %d", rec.Code, http.StatusForbidden)
	}
}
