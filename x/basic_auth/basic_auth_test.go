package basicauth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"codeberg.org/urutau-ltd/aile/v2"
)

func TestMiddlewareRejectsInvalidCredentials(t *testing.T) {
	h := Middleware("", func(user, pass string) bool {
		return user == "admin" && pass == "secret"
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		aile.Status(w, http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("unexpected status: got %d want %d", rec.Code, http.StatusUnauthorized)
	}
	if got := rec.Header().Get("WWW-Authenticate"); got != `Basic realm="restricted"` {
		t.Fatalf("unexpected challenge: got %q", got)
	}
}

func TestMiddlewareAllowsValidCredentials(t *testing.T) {
	h := Middleware("admin", func(user, pass string) bool {
		return user == "admin" && pass == "secret"
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		aile.Status(w, http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.SetBasicAuth("admin", "secret")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("unexpected status: got %d want %d", rec.Code, http.StatusNoContent)
	}
}
