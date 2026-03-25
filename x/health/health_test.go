package health

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"codeberg.org/urutau-ltd/aile/v2"
)

func TestMount(t *testing.T) {
	app, err := aile.New()
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	Mount(app)

	st, err := app.Build(context.Background())
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	tests := []struct {
		path string
		want string
	}{
		{path: "/healthz", want: "ok"},
	}

	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			rec := httptest.NewRecorder()

			st.Server.Handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("unexpected status: got %d want %d", rec.Code, http.StatusOK)
			}
			if rec.Body.String() != tc.want {
				t.Fatalf("unexpected body: got %q want %q", rec.Body.String(), tc.want)
			}
		})
	}
}
