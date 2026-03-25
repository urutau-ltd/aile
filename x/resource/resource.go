package resource

import (
	"errors"
	"net/http"
	"path"
	"reflect"
	"strings"

	"codeberg.org/urutau-ltd/aile/v2"
)

// Collection exposes the conventional handlers used by HTML CRUD resources.
type Collection interface {
	Index(http.ResponseWriter, *http.Request)
	New(http.ResponseWriter, *http.Request)
	Create(http.ResponseWriter, *http.Request)
	Show(http.ResponseWriter, *http.Request)
	Edit(http.ResponseWriter, *http.Request)
	Update(http.ResponseWriter, *http.Request)
	Delete(http.ResponseWriter, *http.Request)
}

// Singleton exposes handlers for resources that should exist only once, such
// as application settings.
type Singleton interface {
	Show(http.ResponseWriter, *http.Request)
	Edit(http.ResponseWriter, *http.Request)
	Update(http.ResponseWriter, *http.Request)
}

// MountCollection registers the conventional CRUD routes for a collection
// resource under basePath.
func MountCollection(a *aile.App, basePath string, h Collection) error {
	basePath, err := normalizeBasePath(basePath)
	if err != nil {
		return err
	}
	if err := validateApp(a); err != nil {
		return err
	}
	if err := validateHandler(h); err != nil {
		return err
	}

	a.GET(basePath, h.Index)
	a.GET(basePath+"/new", h.New)
	a.POST(basePath, h.Create)
	a.GET(basePath+"/{id}", h.Show)
	a.GET(basePath+"/{id}/edit", h.Edit)
	a.POST(basePath+"/{id}/edit", h.Update)
	a.DELETE(basePath+"/{id}", h.Delete)

	return nil
}

// MountSingleton registers the conventional show/edit/update routes for a
// singleton resource under basePath.
func MountSingleton(a *aile.App, basePath string, h Singleton) error {
	basePath, err := normalizeBasePath(basePath)
	if err != nil {
		return err
	}
	if err := validateApp(a); err != nil {
		return err
	}
	if err := validateHandler(h); err != nil {
		return err
	}

	a.GET(basePath, h.Show)
	a.GET(basePath+"/edit", h.Edit)
	a.POST(basePath+"/edit", h.Update)

	return nil
}

func validateApp(a *aile.App) error {
	if a == nil {
		return errors.New("resource: nil app")
	}
	return nil
}

func validateHandler(v any) error {
	if v == nil {
		return errors.New("resource: nil handler")
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Pointer, reflect.Interface, reflect.Map, reflect.Slice, reflect.Func:
		if rv.IsNil() {
			return errors.New("resource: nil handler")
		}
	}

	return nil
}

func normalizeBasePath(basePath string) (string, error) {
	if basePath == "" {
		return "", errors.New("resource: empty base path")
	}
	if !strings.HasPrefix(basePath, "/") {
		return "", errors.New(`resource: base path must start with "/"`)
	}
	if strings.ContainsAny(basePath, "{}") {
		return "", errors.New("resource: base path must be literal")
	}

	clean := path.Clean(basePath)
	if !strings.HasPrefix(clean, "/") {
		clean = "/" + clean
	}
	if clean == "/" {
		return "", errors.New("resource: root base path is not supported")
	}

	return clean, nil
}
