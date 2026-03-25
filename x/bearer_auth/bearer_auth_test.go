package bearerauth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"codeberg.org/urutau-ltd/aile/v2"
)

func TestMiddlewareRejectsInvalidToken(t *testing.T) {
	h := Middleware(func(token string) bool {
		return token == "good-token"
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		aile.Status(w, http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/private", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("unexpected status: got %d want %d", rec.Code, http.StatusUnauthorized)
	}
	if got := rec.Header().Get("WWW-Authenticate"); got != "Bearer" {
		t.Fatalf("unexpected challenge: got %q want %q", got, "Bearer")
	}
}

func TestMiddlewareAllowsValidToken(t *testing.T) {
	h := Middleware(func(token string) bool {
		return token == "good-token"
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		aile.Status(w, http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/private", nil)
	req.Header.Set("Authorization", "Bearer good-token")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("unexpected status: got %d want %d", rec.Code, http.StatusNoContent)
	}
}
