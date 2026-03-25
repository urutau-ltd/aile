package compress

import (
	"bufio"
	"compress/gzip"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"

	"codeberg.org/urutau-ltd/aile/v2"
)

// Config controls gzip compression behavior.
type Config struct {
	Level   int
	MinSize int
}

var gzipPool sync.Pool

// Middleware compresses eligible responses with gzip when the client accepts
// it.
func Middleware(cfg Config) aile.Middleware {
	level := cfg.Level
	if level == 0 {
		level = gzip.DefaultCompression
	}
	minSize := cfg.MinSize

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				next.ServeHTTP(w, r)
				return
			}

			cw := &responseWriter{
				ResponseWriter: w,
				level:          level,
				minSize:        minSize,
			}
			defer cw.Close()
			next.ServeHTTP(cw, r)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
	level       int
	minSize     int
	buf         []byte
	gz          *gzip.Writer
}

func (w *responseWriter) Header() http.Header { return w.ResponseWriter.Header() }

func (w *responseWriter) WriteHeader(status int) {
	if w.wroteHeader {
		return
	}
	w.status = status
	w.wroteHeader = true
}

func (w *responseWriter) Write(p []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	if w.gz == nil && len(w.buf)+len(p) < w.minSize {
		w.buf = append(w.buf, p...)
		return len(p), nil
	}
	if err := w.ensureWriter(); err != nil {
		return 0, err
	}
	return w.gz.Write(p)
}

func (w *responseWriter) ensureWriter() error {
	if w.gz != nil {
		return nil
	}
	w.Header().Del("Content-Length")
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Add("Vary", "Accept-Encoding")
	w.ResponseWriter.WriteHeader(w.status)
	if v := gzipPool.Get(); v != nil {
		gw := v.(*gzip.Writer)
		gw.Reset(w.ResponseWriter)
		w.gz = gw
	} else {
		gw, err := gzip.NewWriterLevel(w.ResponseWriter, w.level)
		if err != nil {
			return err
		}
		w.gz = gw
	}
	if len(w.buf) > 0 {
		if _, err := w.gz.Write(w.buf); err != nil {
			return err
		}
		w.buf = nil
	}
	return nil
}

func (w *responseWriter) Close() error {
	if !w.wroteHeader {
		w.status = http.StatusOK
		w.wroteHeader = true
	}
	if w.gz == nil {
		if len(w.buf) > 0 {
			w.ResponseWriter.Header().Add("Vary", "Accept-Encoding")
			w.ResponseWriter.WriteHeader(w.status)
			_, err := w.ResponseWriter.Write(w.buf)
			w.buf = nil
			return err
		}
		w.ResponseWriter.WriteHeader(w.status)
		return nil
	}
	err := w.gz.Close()
	gzipPool.Put(w.gz)
	w.gz = nil
	return err
}

func (w *responseWriter) Flush() {
	if err := w.ensureWriter(); err == nil {
		_ = w.gz.Flush()
	}
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
	if err := w.ensureWriter(); err != nil {
		return 0, err
	}
	return io.Copy(w.gz, r)
}
