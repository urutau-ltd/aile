package aile_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing/fstest"
	"time"

	aile "codeberg.org/urutau-ltd/aile/v2"
)

func ExampleNew() {
	app, err := aile.New()
	if err != nil {
		panic(err)
	}

	fmt.Println(app != nil)
	// Output: true
}

func ExampleMustNew() {
	app := aile.MustNew()
	fmt.Println(app != nil)
	// Output: true
}

func ExampleDefaultConfig() {
	cfg := aile.DefaultConfig()
	fmt.Println(cfg.Addr)
	// Output: :9001
}

func ExampleWithAddr() {
	app := aile.MustNew(aile.WithAddr(":8080"))
	st, err := app.Build(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(st.Config.Addr)
	// Output: :8080
}

func ExampleWithConfig() {
	app := aile.MustNew(aile.WithConfig(aile.Config{
		Addr:            ":8081",
		ShutdownTimeout: 2 * time.Second,
	}))
	st, err := app.Build(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(st.Config.Addr)
	fmt.Println(st.Config.ShutdownTimeout == 2*time.Second)
	// Output:
	// :8081
	// true
}

func ExampleWithMiddleware() {
	app := aile.MustNew(aile.WithMiddleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-App", "aile")
			next.ServeHTTP(w, r)
		})
	}))
	app.GET("/ping", func(w http.ResponseWriter, r *http.Request) {
		aile.Text(w, http.StatusOK, "ok")
	})

	st, err := app.Build(context.Background())
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	st.Handler.ServeHTTP(rec, req)

	fmt.Println(rec.Header().Get("X-App"))
	// Output: aile
}

func ExampleApp_Use() {
	app := aile.MustNew()
	app.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-App", "aile")
			next.ServeHTTP(w, r)
		})
	})
	app.GET("/ping", func(w http.ResponseWriter, r *http.Request) {
		aile.Text(w, http.StatusOK, "ok")
	})

	st, err := app.Build(context.Background())
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	st.Handler.ServeHTTP(rec, req)

	fmt.Println(rec.Header().Get("X-App"))
	// Output: aile
}

func ExampleApp_POST() {
	app := aile.MustNew()
	app.POST("/articles", func(w http.ResponseWriter, r *http.Request) {
		aile.Status(w, http.StatusCreated)
	})

	st, err := app.Build(context.Background())
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/articles", nil)
	rec := httptest.NewRecorder()
	st.Handler.ServeHTTP(rec, req)

	fmt.Println(rec.Code)
	// Output: 201
}

func ExampleApp_PUT() {
	app := aile.MustNew()
	app.PUT("/articles/1", func(w http.ResponseWriter, r *http.Request) {
		aile.Text(w, http.StatusOK, "updated")
	})

	st, err := app.Build(context.Background())
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest(http.MethodPut, "/articles/1", nil)
	rec := httptest.NewRecorder()
	st.Handler.ServeHTTP(rec, req)

	fmt.Println(rec.Body.String())
	// Output: updated
}

func ExampleApp_PATCH() {
	app := aile.MustNew()
	app.PATCH("/articles/1", func(w http.ResponseWriter, r *http.Request) {
		aile.Text(w, http.StatusOK, "patched")
	})

	st, err := app.Build(context.Background())
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest(http.MethodPatch, "/articles/1", nil)
	rec := httptest.NewRecorder()
	st.Handler.ServeHTTP(rec, req)

	fmt.Println(rec.Body.String())
	// Output: patched
}

func ExampleApp_DELETE() {
	app := aile.MustNew()
	app.DELETE("/articles/1", func(w http.ResponseWriter, r *http.Request) {
		aile.Status(w, http.StatusNoContent)
	})

	st, err := app.Build(context.Background())
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest(http.MethodDelete, "/articles/1", nil)
	rec := httptest.NewRecorder()
	st.Handler.ServeHTTP(rec, req)

	fmt.Println(rec.Code)
	// Output: 204
}

func ExampleApp_HEAD() {
	app := aile.MustNew()
	app.HEAD("/ping", func(w http.ResponseWriter, r *http.Request) {
		aile.Status(w, http.StatusNoContent)
	})

	st, err := app.Build(context.Background())
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest(http.MethodHead, "/ping", nil)
	rec := httptest.NewRecorder()
	st.Handler.ServeHTTP(rec, req)

	fmt.Println(rec.Code, rec.Body.Len())
	// Output: 204 0
}

func ExampleApp_OPTIONS() {
	app := aile.MustNew()
	app.OPTIONS("/ping", func(w http.ResponseWriter, r *http.Request) {
		aile.Status(w, http.StatusNoContent)
	})

	st, err := app.Build(context.Background())
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest(http.MethodOptions, "/ping", nil)
	rec := httptest.NewRecorder()
	st.Handler.ServeHTTP(rec, req)

	fmt.Println(rec.Code)
	// Output: 204
}

func ExampleApp_OnStart() {
	app := aile.MustNew(aile.WithAddr("127.0.0.1:0"))
	app.GET("/ping", func(w http.ResponseWriter, r *http.Request) {
		aile.Text(w, http.StatusOK, "ok")
	})
	app.OnStart(func(ctx context.Context, st *aile.State) error {
		fmt.Println("started")
		shutdownSoon(st)
		return nil
	})

	if err := app.ListenAndServe(); err != nil {
		panic(err)
	}
	// Output: started
}

func ExampleApp_OnShutdown() {
	app := aile.MustNew(aile.WithAddr("127.0.0.1:0"))
	app.GET("/ping", func(w http.ResponseWriter, r *http.Request) {
		aile.Text(w, http.StatusOK, "ok")
	})
	app.OnStart(func(ctx context.Context, st *aile.State) error {
		shutdownSoon(st)
		return nil
	})
	app.OnShutdown(func(ctx context.Context, st *aile.State) error {
		fmt.Println("stopped")
		return nil
	})

	if err := app.ListenAndServe(); err != nil {
		panic(err)
	}
	// Output: stopped
}

func ExampleApp_Set() {
	app := aile.MustNew()
	app.Set("name", "api")

	v, ok := app.Value("name")
	fmt.Println(v, ok)
	// Output: api true
}

func ExampleApp_Value() {
	app := aile.MustNew()
	app.Set("answer", 42)

	v, ok := app.Value("answer")
	fmt.Println(v, ok)
	// Output: 42 true
}

func ExampleApp_Build() {
	app := aile.MustNew()
	app.GET("/ping", func(w http.ResponseWriter, r *http.Request) {
		aile.Text(w, http.StatusOK, "ok")
	})

	st, err := app.Build(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(st.Server != nil && st.Handler != nil && st.Mux != nil)
	// Output: true
}

func ExampleApp_Addr() {
	app := aile.MustNew(aile.WithAddr("127.0.0.1:0"))
	app.GET("/ping", func(w http.ResponseWriter, r *http.Request) {
		aile.Text(w, http.StatusOK, "ok")
	})

	addrCh := make(chan string, 1)
	app.OnStart(func(ctx context.Context, st *aile.State) error {
		addrCh <- app.Addr()
		shutdownSoon(st)
		return nil
	})

	if err := app.ListenAndServe(); err != nil {
		panic(err)
	}

	fmt.Println(<-addrCh != "")
	// Output: true
}

func ExampleApp_Serve() {
	app := aile.MustNew()
	app.GET("/ping", func(w http.ResponseWriter, r *http.Request) {
		aile.Text(w, http.StatusOK, "ok")
	})
	app.OnStart(func(ctx context.Context, st *aile.State) error {
		shutdownSoon(st)
		return nil
	})

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	if err := app.Serve(ln); err != nil {
		panic(err)
	}
}

func ExampleApp_ListenAndServe() {
	app := aile.MustNew(aile.WithAddr("127.0.0.1:0"))
	app.GET("/ping", func(w http.ResponseWriter, r *http.Request) {
		aile.Text(w, http.StatusOK, "ok")
	})
	app.OnStart(func(ctx context.Context, st *aile.State) error {
		shutdownSoon(st)
		return nil
	})

	if err := app.ListenAndServe(); err != nil {
		panic(err)
	}
}

func ExampleApp_Shutdown() {
	app := aile.MustNew()
	app.GET("/ping", func(w http.ResponseWriter, r *http.Request) {
		aile.Text(w, http.StatusOK, "ok")
	})

	ready := make(chan struct{}, 1)
	app.OnStart(func(ctx context.Context, st *aile.State) error {
		ready <- struct{}{}
		return nil
	})

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.Serve(ln)
	}()

	<-ready

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := app.Shutdown(shutdownCtx); err != nil {
		panic(err)
	}

	if err := <-errCh; err != nil {
		panic(err)
	}
}

func ExampleApp_Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := aile.MustNew(aile.WithAddr("127.0.0.1:0"))
	app.GET("/ping", func(w http.ResponseWriter, r *http.Request) {
		aile.Text(w, http.StatusOK, "ok")
	})
	app.OnStart(func(ctx context.Context, st *aile.State) error {
		go func() {
			time.Sleep(10 * time.Millisecond)
			cancel()
		}()
		return nil
	})

	if err := app.Run(ctx); err != nil {
		panic(err)
	}
}

func ExampleStaticHandler() {
	h, err := aile.StaticHandler("/assets/", fstest.MapFS{
		"app.css": &fstest.MapFile{Data: []byte("body{color:black}")},
	})
	if err != nil {
		panic(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/assets/app.css", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	fmt.Println(rec.Code, rec.Body.String())
	// Output: 200 body{color:black}
}

func ExampleStatus() {
	rec := httptest.NewRecorder()
	aile.Status(rec, http.StatusNoContent)

	fmt.Println(rec.Code)
	// Output: 204
}

func ExampleText() {
	rec := httptest.NewRecorder()
	aile.Text(rec, http.StatusOK, "hello")

	fmt.Println(rec.Header().Get("Content-Type"))
	fmt.Println(rec.Body.String())
	// Output:
	// text/plain; charset=utf-8
	// hello
}

func ExampleError() {
	rec := httptest.NewRecorder()
	aile.Error(rec, http.StatusBadRequest, "bad request")

	fmt.Println(rec.Code)
	fmt.Print(strings.TrimSpace(rec.Body.String()))
	// Output:
	// 400
	// bad request
}

func ExampleReleaseTag() {
	fmt.Println(aile.ReleaseTag)
	// Output: v2.1.0
}

func shutdownSoon(st *aile.State) {
	go func() {
		time.Sleep(10 * time.Millisecond)
		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_ = st.Server.Shutdown(shutdownCtx)
	}()
}

func ExampleApp_Serve_request() {
	app := aile.MustNew()
	app.GET("/ping", func(w http.ResponseWriter, r *http.Request) {
		aile.Text(w, http.StatusOK, "ok")
	})

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	done := make(chan error, 1)
	app.OnStart(func(ctx context.Context, st *aile.State) error {
		go func() {
			time.Sleep(10 * time.Millisecond)
			resp, err := http.Get("http://" + ln.Addr().String() + "/ping")
			if err != nil {
				done <- err
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				done <- err
				return
			}

			fmt.Println(resp.StatusCode, string(body))

			shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			done <- st.Server.Shutdown(shutdownCtx)
		}()
		return nil
	})

	if err := app.Serve(ln); err != nil {
		panic(err)
	}
	if err := <-done; err != nil {
		panic(err)
	}
	// Output: 200 ok
}
