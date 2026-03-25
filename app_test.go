package aile

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"runtime"
	"testing"
)

func TestNewDefaults(t *testing.T) {
	app, err := New()
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	if app.config.Addr == "" {
		t.Fatal("expected default addr to be set")
	}
	if app.values == nil {
		t.Fatal("expected values map to be initialized")
	}
}

func TestBuildRegistersLiteralServeMuxPatterns(t *testing.T) {
	app, err := New()

	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	app.GET("/hello", func(w http.ResponseWriter, r *http.Request) {
		Text(w, http.StatusOK, "ok")
	})

	st, err := app.Build(context.Background())
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	rec := debugServeMuxRequest(t, st.Mux, http.MethodGet, "/hello")

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: got %d want %d", rec.Code, http.StatusOK)
	}

	if rec.Body.String() != "ok" {
		t.Fatalf("unexpected body: got %q want %q", rec.Body.String(), "ok")
	}

}

func TestGetHelperRegistersMethodPattern(t *testing.T) {
	app := MustNew()

	app.GET("/hello", func(w http.ResponseWriter, r *http.Request) {
		Text(w, http.StatusOK, "ok")
	})

	st, err := app.Build(context.Background())
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	rec := httptest.NewRecorder()
	st.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: got %d want %d", rec.Code, http.StatusOK)
	}
}

func TestBuildAppliesMiddlewareInExpectedOrder(t *testing.T) {
	app := MustNew()

	var got []string

	app.Use(
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

	app.GET("/x", func(w http.ResponseWriter, r *http.Request) {
		got = append(got, "handler")
		Status(w, http.StatusNoContent)
	})

	st, err := app.Build(context.Background())
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	rec := httptest.NewRecorder()
	st.Handler.ServeHTTP(rec, req)

	want := []string{"mw1-before", "mw2-before", "handler", "mw2-after", "mw1-after"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected middleware order\n got: %#v\nwant: %#v", got, want)
	}
}

func TestBuildAppliesConfigDefaults(t *testing.T) {
	app := MustNew(WithConfig(Config{}))

	st, err := app.Build(context.Background())
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	def := DefaultConfig()
	if st.Config.Addr != def.Addr {
		t.Fatalf("unexpected addr: got %q want %q", st.Config.Addr, def.Addr)
	}
	if st.Config.ReadTimeout != def.ReadTimeout {
		t.Fatalf("unexpected ReadTimeout: got %v want %v", st.Config.ReadTimeout, def.ReadTimeout)
	}
	if st.Config.ShutdownTimeout != def.ShutdownTimeout {
		t.Fatalf("unexpected ShutdownTimeout: got %v want %v", st.Config.ShutdownTimeout, def.ShutdownTimeout)
	}
}

func TestBuildCopiesValues(t *testing.T) {
	app := MustNew()
	app.Set("name", "aile")

	st, err := app.Build(context.Background())
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	st.Values["name"] = "other"

	got, ok := app.Value("name")
	if !ok {
		t.Fatal("expected value to exist")
	}
	if got != "aile" {
		t.Fatalf("unexpected app value after state mutation: got %#v want %#v", got, "aile")
	}
}

func TestBuildFailsOnEmptyPath(t *testing.T) {
	app := MustNew()
	app.GET("", func(w http.ResponseWriter, r *http.Request) {})

	_, err := app.Build(context.Background())
	if err == nil {
		t.Fatal("expected build error for empty route path")
	}
}

func TestBuildFailsOnNilHandler(t *testing.T) {
	app := MustNew()
	app.GET("/x", nil)

	_, err := app.Build(context.Background())
	if err == nil {
		t.Fatal("expected build error for nil handler")
	}
}

func TestSetValueAndValue(t *testing.T) {
	app := MustNew()
	app.Set("answer", 42)

	got, ok := app.Value("answer")
	if !ok {
		t.Fatal("expected value to exist")
	}
	if got != 42 {
		t.Fatalf("unexpected value: got %#v want %#v", got, 42)
	}
}

func TestServeMuxMethodPatternSmoke(t *testing.T) {
	t.Logf("go version: %s", runtime.Version())
	t.Logf("GOOS=%s GOARCH=%s", runtime.GOOS, runtime.GOARCH)
	t.Logf("GODEBUG=%q", os.Getenv("GODEBUG"))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	req := httptest.NewRequest(http.MethodGet, "/hello", nil)

	h, pattern := mux.Handler(req)
	t.Logf("resolved pattern=%q handler=%T", pattern, h)

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	t.Logf("status=%d body=%q", rec.Code, rec.Body.String())

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: got %d want %d", rec.Code, http.StatusOK)
	}
}
