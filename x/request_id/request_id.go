package requestid

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"codeberg.org/urutau-ltd/aile"
)

type contextKey struct{}

type Config struct {
	Header    string
	Generator func() string
}

var fallbackCounter uint64

func Middleware(cfg Config) aile.Middleware {
	header := cfg.Header
	if header == "" {
		header = "X-Request-ID"
	}
	gen := cfg.Generator
	if gen == nil {
		gen = defaultGenerator
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.Header.Get(header)
			if id == "" {
				id = gen()
			}
			ctx := context.WithValue(r.Context(), contextKey{}, id)
			w.Header().Set(header, id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func FromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(contextKey{}).(string)
	return v, ok && v != ""
}

func defaultGenerator() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err == nil {
		return hex.EncodeToString(b[:])
	}
	n := atomic.AddUint64(&fallbackCounter, 1)
	return strconv.FormatInt(time.Now().UnixNano(), 16) + "-" + strconv.FormatUint(n, 16)
}
