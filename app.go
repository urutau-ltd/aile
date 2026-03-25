package aile

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
)

// Middleware is the standard [net/http] middleware shape.
type Middleware func(http.Handler) http.Handler

// HookFunc runs during startup or shutdown.
type HookFunc func(context.Context, *State) error

type routeDef struct {
	method  string
	path    string
	handler http.HandlerFunc
}

// App is the main application container.
//
// App is safe for concurrent reads after construction, but configuration
// methods such as route registration, [App.Use], [App.Set], [App.OnStart] and
// [App.OnShutdown] should be called before starting the server itself.
//
// Configure the app before starting the server. Mutating routes, middleware,
// values or hooks while serving is NOT supported.
type App struct {
	mu sync.RWMutex

	config Config

	routes []routeDef
	mws    []Middleware
	values map[string]any

	onStart    []HookFunc
	onShutdown []HookFunc

	running *State
	ln      net.Listener
}

// New constructs a new [App] and applies the provided options.
func New(opts ...Option) (*App, error) {
	a := &App{
		config: DefaultConfig(),
		values: make(map[string]any),
	}

	for _, opt := range opts {
		if err := opt(a); err != nil {
			return nil, err
		}
	}

	return a, nil
}

// MustNew is like [New] but panics if construction fails.
func MustNew(opts ...Option) *App {
	a, err := New(opts...)
	if err != nil {
		panic(err)
	}
	return a
}

// Use appends middleware to the app in registration order.
func (a *App) Use(mw ...Middleware) {
	a.mws = append(a.mws, mw...)
}

func (a *App) route(method, path string, h http.HandlerFunc) {
	a.routes = append(a.routes, routeDef{
		method:  method,
		path:    path,
		handler: h,
	})
}

// GET registers a GET route in an App.
func (a *App) GET(pattern string, h http.HandlerFunc) {
	a.route(http.MethodGet, pattern, h)
}

// POST registers a POST route in an App.
func (a *App) POST(pattern string, h http.HandlerFunc) {
	a.route(http.MethodPost, pattern, h)
}

// PUT registers a PUT route in an App.
func (a *App) PUT(pattern string, h http.HandlerFunc) {
	a.route(http.MethodPut, pattern, h)
}

// PATCH registers a PATCH route in an App.
func (a *App) PATCH(pattern string, h http.HandlerFunc) {
	a.route(http.MethodPatch, pattern, h)
}

// DELETE registers a DELETE route in an App.
func (a *App) DELETE(pattern string, h http.HandlerFunc) {
	a.route(http.MethodDelete, pattern, h)
}

// HEAD registers a HEAD route in an App.
func (a *App) HEAD(pattern string, h http.HandlerFunc) {
	a.route(http.MethodHead, pattern, h)
}

// OPTIONS registers an OPTIONS route in an App.
func (a *App) OPTIONS(pattern string, h http.HandlerFunc) {
	a.route(http.MethodOptions, pattern, h)
}

// Use OnStart to append a startup hook.
func (a *App) OnStart(fn HookFunc) {
	a.onStart = append(a.onStart, fn)
}

// Use OnShutdown to append a shutdown hook.
func (a *App) OnShutdown(fn HookFunc) {
	a.onShutdown = append(a.onShutdown, fn)
}

// Set stores a named value on the App.
func (a *App) Set(name string, value any) {
	if a.values == nil {
		a.values = make(map[string]any)
	}
	a.values[name] = value
}

// Value searches for a named value from the app.
func (a *App) Value(name string) (any, bool) {
	v, ok := a.values[name]
	return v, ok
}

// Build constructs the app into a [State] without starting the server.
func (a *App) Build(ctx context.Context) (*State, error) {
	_ = ctx

	cfg := a.config
	applyConfigDefaults(&cfg)

	mux := http.NewServeMux()

	for _, rt := range a.routes {
		if rt.method == "" {
			return nil, errors.New("build: empty route method")
		}

		if rt.path == "" {
			return nil, fmt.Errorf("build: empty route path for method %q", rt.method)
		}

		if rt.handler == nil {
			return nil, fmt.Errorf("build: nil handler for route %q", rt.method+" "+rt.path)
		}

		mux.Handle(rt.method+" "+rt.path, rt.handler)
	}

	var handler http.Handler = mux
	handler = chain(handler, a.mws...)

	server := &http.Server{
		Addr:              cfg.Addr,
		Handler:           handler,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		MaxHeaderBytes:    cfg.MaxHeaderBytes,
	}

	st := &State{
		Config:  cfg,
		Mux:     mux,
		Handler: handler,
		Server:  server,
		Values:  cloneMap(a.values),
	}

	return st, nil
}

func chain(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

func cloneMap(in map[string]any) map[string]any {
	if in == nil {
		return nil
	}
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func (a *App) setRunning(st *State, ln net.Listener) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.running = st
	a.ln = ln
}

func (a *App) clearRunning() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.running = nil
	a.ln = nil
}
