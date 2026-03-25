package secureheaders

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"

	"codeberg.org/urutau-ltd/aile/v2"
)

func TestMiddlewareSetsConfiguredHeaders(t *testing.T) {
	h := Middleware(Config{
		ContentTypeNosniff:      true,
		FrameDeny:               true,
		ReferrerPolicy:          "same-origin",
		ContentSecurityPolicy:   "default-src 'self'",
		PermissionsPolicy:       "camera=()",
		CrossOriginOpenerPolicy: "same-origin",
		HSTSMaxAge:              60,
		HSTSIncludeSubdomains:   true,
		HSTSPreload:             true,
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		aile.Status(w, http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.TLS = &tls.ConnectionState{}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	tests := map[string]string{
		"X-Content-Type-Options":     "nosniff",
		"X-Frame-Options":            "DENY",
		"Referrer-Policy":            "same-origin",
		"Content-Security-Policy":    "default-src 'self'",
		"Permissions-Policy":         "camera=()",
		"Cross-Origin-Opener-Policy": "same-origin",
		"Strict-Transport-Security":  "max-age=60; includeSubDomains; preload",
	}

	for header, want := range tests {
		if got := rec.Header().Get(header); got != want {
			t.Fatalf("unexpected %s: got %q want %q", header, got, want)
		}
	}
}

func TestMiddlewareSkipsHSTSWithoutTLS(t *testing.T) {
	h := Middleware(Config{HSTSMaxAge: 60})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		aile.Status(w, http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if got := rec.Header().Get("Strict-Transport-Security"); got != "" {
		t.Fatalf("unexpected HSTS header: got %q want empty", got)
	}
}
