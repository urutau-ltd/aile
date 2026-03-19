package aile

// The Option function configures an [aile/App] during construction
type Option func(*App) error

// WithAddr sets the server listen address, overriding the default
// one which is ":9001"
func WithAddr(addr string) Option {
	return func(a *App) error {
		a.config.Addr = addr
		return nil
	}
}

// WithConfig replaces the full app configuration.
// Zero values are later filled using [aile/DefaultConfig] during Build/Run.
func WithConfig(cfg Config) Option {
	return func(a *App) error {
		a.config = cfg
		return nil
	}
}

// WithMiddleware appends middleware which are [http/HandlerFunction] during
// construction
func WithMiddleware(mw ...Middleware) Option {
	return func(a *App) error {
		a.Use(mw...)
		return nil
	}
}
