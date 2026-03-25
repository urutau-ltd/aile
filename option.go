package aile

// Option configures an [App] during construction.
type Option func(*App) error

// WithAddr sets the server listen address, overriding the default ":9001".
func WithAddr(addr string) Option {
	return func(a *App) error {
		a.config.Addr = addr
		return nil
	}
}

// WithConfig replaces the full app configuration.
// Zero values are later filled using [DefaultConfig] during Build or Run.
func WithConfig(cfg Config) Option {
	return func(a *App) error {
		a.config = cfg
		return nil
	}
}

// WithMiddleware appends middleware during construction.
func WithMiddleware(mw ...Middleware) Option {
	return func(a *App) error {
		a.Use(mw...)
		return nil
	}
}
