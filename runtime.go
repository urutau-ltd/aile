package aile

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// Addr returns the current listener address for a running app.
//
// If the app is not running, Addr returns the empty string.
func (a *App) Addr() string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.ln == nil {
		return ""
	}

	return a.ln.Addr().String()
}

// Serve builds the app and serves it on the provided listener.
func (a *App) Serve(ln net.Listener) error {
	st, err := a.Build(context.Background())

	if err != nil {
		return err
	}

	a.setRunning(st, ln)
	defer a.clearRunning()

	if err := runHooks(context.Background(), st, a.onStart); err != nil {
		_ = ln.Close()
		return err
	}

	err = st.Server.Serve(ln)
	if errors.Is(err, http.ErrServerClosed) {
		err = nil
	}

	hookErr := runHooks(context.Background(), st, a.onShutdown)
	return errors.Join(err, hookErr)
}

// ListenAndServe listens on the configured address and then serves the app.
func (a *App) ListenAndServe() error {
	st, err := a.Build(context.Background())
	if err != nil {
		return err
	}

	ln, err := net.Listen("tcp", st.Config.Addr)
	if err != nil {
		return err
	}

	return a.Serve(ln)
}

// Shutdown gracefully stops a running app.
// If the app is not running, it returns nil.
func (a *App) Shutdown(ctx context.Context) error {
	a.mu.RLock()
	st := a.running
	a.mu.RUnlock()

	if st == nil || st.Server == nil {
		return nil
	}

	return st.Server.Shutdown(ctx)
}

// Run listens on the configured address, serves requests, handles shutdown
// signals and performs graceful shutdown.
func (a *App) Run(ctx context.Context) error {
	st, err := a.Build(ctx)

	if err != nil {
		return err
	}

	ln, err := net.Listen("tcp", st.Config.Addr)
	if err != nil {
		return err
	}

	a.setRunning(st, ln)
	defer a.clearRunning()

	if err := runHooks(ctx, st, a.onStart); err != nil {
		_ = ln.Close()
		return err
	}

	runCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		serveErr := st.Server.Serve(ln)
		if errors.Is(serveErr, http.ErrServerClosed) {
			serveErr = nil
		}
		errCh <- serveErr
	}()

	var runErr error

	select {
	case runErr = <-errCh:
	case <-runCtx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), st.Config.ShutdownTimeout)
		defer cancel()

		if err := st.Server.Shutdown(shutdownCtx); err != nil {
			runErr = err
		} else {
			runErr = <-errCh
		}
	}

	hookErr := runHooks(context.Background(), st, a.onShutdown)
	return errors.Join(runErr, hookErr)
}

func runHooks(ctx context.Context, st *State, hooks []HookFunc) error {
	for _, fn := range hooks {
		if err := fn(ctx, st); err != nil {
			return err
		}
	}
	return nil
}
