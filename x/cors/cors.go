package cors

import (
	"net/http"
	"strconv"
	"strings"

	"codeberg.org/urutau-ltd/aile/v2"
)

// Config controls CORS headers emitted by the middleware.
type Config struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           int
}

// Middleware applies a simple CORS policy to incoming requests.
func Middleware(cfg Config) aile.Middleware {
	allowMethods := strings.Join(defaultSlice(cfg.AllowMethods, []string{
		http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch,
		http.MethodDelete, http.MethodHead, http.MethodOptions,
	}), ", ")
	allowHeaders := strings.Join(cfg.AllowHeaders, ", ")
	exposeHeaders := strings.Join(cfg.ExposeHeaders, ", ")
	origins := defaultSlice(cfg.AllowOrigins, []string{"*"})

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" && allowedOrigin(origin, origins) {
				if len(origins) == 1 && origins[0] == "*" {
					w.Header().Set("Access-Control-Allow-Origin", "*")
				} else {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Add("Vary", "Origin")
				}
				if cfg.AllowCredentials {
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}
				if exposeHeaders != "" {
					w.Header().Set("Access-Control-Expose-Headers", exposeHeaders)
				}
			}

			if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
				if allowedOrigin(origin, origins) {
					w.Header().Set("Access-Control-Allow-Methods", allowMethods)
					if allowHeaders != "" {
						w.Header().Set("Access-Control-Allow-Headers", allowHeaders)
					}
					if cfg.MaxAge > 0 {
						w.Header().Set("Access-Control-Max-Age", strconv.Itoa(cfg.MaxAge))
					}
				}
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func allowedOrigin(origin string, allowed []string) bool {
	if origin == "" {
		return false
	}
	for _, candidate := range allowed {
		if candidate == "*" || candidate == origin {
			return true
		}
	}
	return false
}

func defaultSlice(got, fallback []string) []string {
	if len(got) == 0 {
		return fallback
	}
	return got
}
