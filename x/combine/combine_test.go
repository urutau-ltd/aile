package combine

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"codeberg.org/urutau-ltd/aile/v2"
)

func TestMiddlewareAppliesWrappedMiddlewaresInOrder(t *testing.T) {
	var got []string

	stack := Middleware(
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				got = append(got, "mw1-before")
				next.ServeHTTP(w, r)
				got = append(got, "mw1-after")
			})
		},
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				got = append(got, "mw2-before")
				next.ServeHTTP(w, r)
				got = append(got, "mw2-after")
			})
		},
	)

	h := stack(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = append(got, "handler")
		aile.Status(w, http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	want := []string{"mw1-before", "mw2-before", "handler", "mw2-after", "mw1-after"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected middleware order\n got: %#v\nwant: %#v", got, want)
	}
}
