package aile

import "time"

const (
	// The HTTP server timeout for read operations duration in seconds.
	// if not set in the application configuration this constant will be used
	// as the fallback value.
	//
	// This is used in the ReadTimeout and ReadHeaderTimeout [Config] fields.
	READ_TIMEOUT  time.Duration = 5 * time.Second
	
	// The HTTP server timeout duration for write operations duration in
	// seconds.
	// if not set in the application configuration this constant will be
	// used as the fallback value.
	//
	// This is used in the WriteTimeout and ShutdownTimeout [Config] fields.
	WRITE_TIMEOUT time.Duration = 10 * time.Second
	LONG_TIMEOUT  time.Duration = 60 * time.Second
)

// Config controls the HTTP server runtime.
type Config struct {
	// An HTTP server address to serve the aile [App].
	Addr              string
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration
	MaxHeaderBytes    int
}

// DefaultConfig returns the default runtime configuration.
func DefaultConfig() Config {
	return Config{
		Addr:              ":9001",
		ReadTimeout:       READ_TIMEOUT,
		ReadHeaderTimeout: READ_TIMEOUT,
		WriteTimeout:      WRITE_TIMEOUT,
		IdleTimeout:       LONG_TIMEOUT,
		ShutdownTimeout:   WRITE_TIMEOUT,
		MaxHeaderBytes:    1 << 20, // 1 MiB
	}
}

func applyConfigDefaults(cfg *Config) {
	def := DefaultConfig()

	if cfg.Addr == "" {
		cfg.Addr = def.Addr
	}
	if cfg.ReadTimeout == 0 {
		cfg.ReadTimeout = def.ReadTimeout
	}
	if cfg.ReadHeaderTimeout == 0 {
		cfg.ReadHeaderTimeout = def.ReadHeaderTimeout
	}
	if cfg.WriteTimeout == 0 {
		cfg.WriteTimeout = def.WriteTimeout
	}
	if cfg.IdleTimeout == 0 {
		cfg.IdleTimeout = def.IdleTimeout
	}
	if cfg.ShutdownTimeout == 0 {
		cfg.ShutdownTimeout = def.ShutdownTimeout
	}
	if cfg.MaxHeaderBytes == 0 {
		cfg.MaxHeaderBytes = def.MaxHeaderBytes
	}
}
