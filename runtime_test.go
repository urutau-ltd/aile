package aile

import (
	"context"
	"errors"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestShutdownBeforeBuildIsNoop(t *testing.T) {
	app := MustNew()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		t.Fatalf("Shutdown returned error: %v", err)
	}
}

func TestServeCallsStartAndShutdownHooks(t *testing.T) {
	app := MustNew(WithAddr("127.0.0.1:0"))
	app.GET("/ping", func(w http.ResponseWriter, r *http.Request) {
		Text(w, http.StatusOK, "ok")
	})

	started := make(chan struct{}, 1)
	stopped := make(chan struct{}, 1)

	app.OnStart(func(ctx context.Context, st *State) error {
		started <- struct{}{}
		go func() {
			time.Sleep(50 * time.Millisecond)
			shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			_ = st.Server.Shutdown(shutdownCtx)
		}()
		return nil
	})

	app.OnShutdown(func(ctx context.Context, st *State) error {
		stopped <- struct{}{}
		return nil
	})

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("net.Listen returned error: %v", err)
	}
	defer ln.Close()

	if err := app.Serve(ln); err != nil {
		t.Fatalf("Serve returned error: %v", err)
	}

	select {
	case <-started:
	default:
		t.Fatal("expected start hook to run")
	}

	select {
	case <-stopped:
	default:
		t.Fatal("expected shutdown hook to run")
	}
}

func TestServeSetsRealAddr(t *testing.T) {
	app := MustNew(WithAddr("127.0.0.1:0"))
	app.GET("/ping", func(w http.ResponseWriter, r *http.Request) {
		Text(w, http.StatusOK, "ok")
	})

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("net.Listen returned error: %v", err)
	}
	defer ln.Close()

	ready := make(chan string, 1)
	app.OnStart(func(ctx context.Context, st *State) error {
		ready <- app.Addr()
		go func() {
			time.Sleep(50 * time.Millisecond)
			shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			_ = st.Server.Shutdown(shutdownCtx)
		}()
		return nil
	})

	if err := app.Serve(ln); err != nil {
		t.Fatalf("Serve returned error: %v", err)
	}

	addr := <-ready
	if addr == "" {
		t.Fatal("expected real runtime addr to be available")
	}
}

func TestServeReturnsShutdownHookError(t *testing.T) {
	app := MustNew(WithAddr("127.0.0.1:0"))
	app.GET("/ping", func(w http.ResponseWriter, r *http.Request) {
		Text(w, http.StatusOK, "ok")
	})

	want := errors.New("shutdown hook failed")

	app.OnStart(func(ctx context.Context, st *State) error {
		go func() {
			time.Sleep(50 * time.Millisecond)
			shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			_ = st.Server.Shutdown(shutdownCtx)
		}()
		return nil
	})

	app.OnShutdown(func(ctx context.Context, st *State) error {
		return want
	})

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("net.Listen returned error: %v", err)
	}
	defer ln.Close()

	err = app.Serve(ln)
	if !errors.Is(err, want) {
		t.Fatalf("Serve error = %v, want error containing %v", err, want)
	}
}

func TestRunReturnsShutdownHookError(t *testing.T) {
	app := MustNew(WithAddr("127.0.0.1:0"))
	app.GET("/ping", func(w http.ResponseWriter, r *http.Request) {
		Text(w, http.StatusOK, "ok")
	})

	want := errors.New("shutdown hook failed")

	app.OnStart(func(ctx context.Context, st *State) error {
		go func() {
			time.Sleep(50 * time.Millisecond)
			shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			_ = st.Server.Shutdown(shutdownCtx)
		}()
		return nil
	})

	app.OnShutdown(func(ctx context.Context, st *State) error {
		return want
	})

	err := app.Run(context.Background())
	if !errors.Is(err, want) {
		t.Fatalf("Run error = %v, want error containing %v", err, want)
	}
}
