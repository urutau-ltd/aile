package logger

import (
	"bufio"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"time"

	"codeberg.org/urutau-ltd/aile/v2"
)

type responseWriter struct {
	http.ResponseWriter
	status      int
	bytes       int64
	wroteHeader bool
}

// Middleware logs basic request and response information using slog.
func Middleware(l *slog.Logger) aile.Middleware {
	if l == nil {
		l = slog.Default()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				start := time.Now()

				ww := &responseWriter{
					ResponseWriter: w,
					status:         http.StatusOK,
				}

				next.ServeHTTP(ww, r)

				l.Info("http request",
					"method", r.Method,
					"path", r.URL.Path,
					"status", ww.status,
					"bytes", ww.bytes,
					"duration", time.Since(start),
					"remote_addr", r.RemoteAddr,
				)
			})
	}
}

func (w *responseWriter) WriteHeader(status int) {
	if w.wroteHeader {
		return
	}

	w.status = status
	w.wroteHeader = true
	w.ResponseWriter.WriteHeader(status)
}

func (w *responseWriter) Write(p []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err := w.ResponseWriter.Write(p)
	w.bytes += int64(n)
	return n, err
}

func (w *responseWriter) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (w *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := w.ResponseWriter.(http.Hijacker)

	if !ok {
		return nil, nil, errors.New("response writer does not support hijacking")
	}

	return h.Hijack()
}

func (w *responseWriter) Push(target string, opts *http.PushOptions) error {
	p, ok := w.ResponseWriter.(http.Pusher)
	if !ok {
		return http.ErrNotSupported
	}

	return p.Push(target, opts)
}

func (w *responseWriter) ReadFrom(r io.Reader) (int64, error) {
	if rf, ok := w.ResponseWriter.(io.ReaderFrom); ok {
		if !w.wroteHeader {
			w.WriteHeader(http.StatusOK)
		}

		n, err := rf.ReadFrom(r)
		w.bytes += n
		return n, err
	}

	return io.Copy(w, r)
}
