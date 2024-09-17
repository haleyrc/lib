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
func (router Router) GET(path string, h http.Handler) {
	router.mux.Handle(fmt.Sprintf("GET %s", path), h)
}

// POST registers a route at the provided path that responds to POST requests.
func (router Router) POST(path string, h http.Handler) {
	router.mux.Handle(fmt.Sprintf("POST %s", path), h)
}

// ServeHTTP implements the http.Handler interface.
func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.mux.ServeHTTP(w, r)
}
