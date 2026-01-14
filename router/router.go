// Package router provides an http.Handler implementation with a nice API for
// registering routes.
package router

import (
	"fmt"
	"net/http"
)

// Router is a wrapper around an http.ServeMux that provides a slightly nicer
// API for registering routes.
type Router struct {
	mux *http.ServeMux
}

// New creates a new Router.
func New() Router {
	return Router{mux: http.NewServeMux()}
}

// GET registers a route at the provided path that responds to GET requests.
func (router Router) GET(path string, h http.HandlerFunc) {
	router.mux.HandleFunc(fmt.Sprintf("GET %s", path), h)
}

// POST registers a route at the provided path that responds to POST requests.
func (router Router) POST(path string, h http.HandlerFunc) {
	router.mux.HandleFunc(fmt.Sprintf("POST %s", path), h)
}

// ServeFiles serves static files located in `dir` and accessed via requests
// whose path starts with `prefix`.
func (router Router) ServeFiles(prefix, dir string) {
	route := fmt.Sprintf("GET %s", prefix)
	fs := http.FileServer(http.Dir(dir))
	router.mux.Handle(route, http.StripPrefix(prefix, fs))
}

// ServeHTTP implements the http.Handler interface.
func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.mux.ServeHTTP(w, r)
}
