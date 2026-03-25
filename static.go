package aile

import (
	"errors"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

// StaticHandler returns an [http.Handler] that serves files from fsys below the
// given URL prefix using the standard library file server.
//
// The prefix must start with "/" and may be passed with or without a trailing
// slash. When the prefix is "/", the returned handler is the plain
// [http.FileServerFS] handler.
func StaticHandler(prefix string, fsys fs.FS) (http.Handler, error) {
	_, routePrefix, _, err := normalizeStaticPrefix(prefix)
	if err != nil {
		return nil, err
	}
	if fsys == nil {
		return nil, errors.New("static: nil fs")
	}
	if routePrefix == "/" {
		return http.FileServerFS(fsys), nil
	}
	return http.StripPrefix(routePrefix, http.FileServerFS(fsys)), nil
}

// Static mounts a file server under prefix.
//
// It wraps [StaticHandler] and registers a GET subtree route on the app. If
// prefix is passed without a trailing slash, it also registers a permanent
// redirect from prefix to prefix+"/".
func (a *App) Static(prefix string, fsys fs.FS) error {
	barePrefix, routePrefix, redirectBare, err := normalizeStaticPrefix(prefix)
	if err != nil {
		return err
	}

	h, err := StaticHandler(prefix, fsys)
	if err != nil {
		return err
	}

	if redirectBare {
		a.GET(barePrefix, func(w http.ResponseWriter, r *http.Request) {
			redirectURL := *r.URL
			redirectURL.Path = routePrefix
			http.Redirect(w, r, redirectURL.String(), http.StatusPermanentRedirect)
		})
	}

	a.GET(routePrefix, h.ServeHTTP)
	return nil
}

func normalizeStaticPrefix(prefix string) (barePrefix, routePrefix string, redirectBare bool, err error) {
	if prefix == "" {
		return "", "", false, errors.New(`static: empty prefix`)
	}
	if !strings.HasPrefix(prefix, "/") {
		return "", "", false, errors.New(`static: prefix must start with "/"`)
	}
	if prefix == "/" {
		return "/", "/", false, nil
	}

	barePrefix = path.Clean(prefix)
	if !strings.HasPrefix(barePrefix, "/") {
		barePrefix = "/" + barePrefix
	}

	routePrefix = barePrefix
	if !strings.HasSuffix(routePrefix, "/") {
		routePrefix += "/"
	}

	return barePrefix, routePrefix, !strings.HasSuffix(prefix, "/"), nil
}
