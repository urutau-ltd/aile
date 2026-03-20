package secureheaders

import (
	"net/http"
	"strconv"

	"codeberg.org/urutau-ltd/aile"
)

type Config struct {
	ContentTypeNosniff      bool
	FrameDeny               bool
	ReferrerPolicy          string
	ContentSecurityPolicy   string
	PermissionsPolicy       string
	CrossOriginOpenerPolicy string
	HSTSMaxAge              int
	HSTSIncludeSubdomains   bool
	HSTSPreload             bool
}

func Middleware(cfg Config) aile.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cfg.ContentTypeNosniff {
				w.Header().Set("X-Content-Type-Options", "nosniff")
			}
			if cfg.FrameDeny {
				w.Header().Set("X-Frame-Options", "DENY")
			}
			if cfg.ReferrerPolicy != "" {
				w.Header().Set("Referrer-Policy", cfg.ReferrerPolicy)
			}
			if cfg.ContentSecurityPolicy != "" {
				w.Header().Set("Content-Security-Policy", cfg.ContentSecurityPolicy)
			}
			if cfg.PermissionsPolicy != "" {
				w.Header().Set("Permissions-Policy", cfg.PermissionsPolicy)
			}
			if cfg.CrossOriginOpenerPolicy != "" {
				w.Header().Set("Cross-Origin-Opener-Policy", cfg.CrossOriginOpenerPolicy)
			}
			if cfg.HSTSMaxAge > 0 && r.TLS != nil {
				value := "max-age=" + strconv.Itoa(cfg.HSTSMaxAge)
				if cfg.HSTSIncludeSubdomains {
					value += "; includeSubDomains"
				}
				if cfg.HSTSPreload {
					value += "; preload"
				}
				w.Header().Set("Strict-Transport-Security", value)
			}
			next.ServeHTTP(w, r)
		})
	}
}
