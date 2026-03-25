package iprestriction

import (
	"net"
	"net/http"
	"strings"

	"codeberg.org/urutau-ltd/aile/v2"
)

// Config controls allow and deny lists for client IP filtering.
type Config struct {
	Allow      []*net.IPNet
	Deny       []*net.IPNet
	TrustProxy bool
}

// Middleware filters requests by client IP address or network.
func Middleware(cfg Config) aile.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, ok := clientIP(r, cfg.TrustProxy)
			if !ok {
				aile.Error(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
				return
			}
			if containsIP(cfg.Deny, ip) {
				aile.Error(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
				return
			}
			if len(cfg.Allow) > 0 && !containsIP(cfg.Allow, ip) {
				aile.Error(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func containsIP(networks []*net.IPNet, ip net.IP) bool {
	for _, n := range networks {
		if n != nil && n.Contains(ip) {
			return true
		}
	}
	return false
}

func clientIP(r *http.Request, trustProxy bool) (net.IP, bool) {
	if trustProxy {
		xff := r.Header.Get("X-Forwarded-For")
		if xff != "" {
			first := strings.TrimSpace(strings.Split(xff, ",")[0])
			if ip := net.ParseIP(first); ip != nil {
				return ip, true
			}
		}
		xri := strings.TrimSpace(r.Header.Get("X-Real-IP"))
		if ip := net.ParseIP(xri); ip != nil {
			return ip, true
		}
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		if ip := net.ParseIP(r.RemoteAddr); ip != nil {
			return ip, true
		}
		return nil, false
	}
	ip := net.ParseIP(host)
	return ip, ip != nil
}
