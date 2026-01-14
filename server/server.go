// Package server provides an HTTP server with sane defaults.
package server

import (
	"context"
	"net"
	"net/http"
	"time"
)

// Server is a wrapper around an http.Server that attempts to shutdown
// gracefully. Unlike http.Server, this implementation does not return an error
// if the server is shutdown cleanly.
type Server struct {
	s *http.Server
}

// New returns a Server that listens on the given port on any interface
// (0.0.0.0:PORT). Read, write, and idle timeouts are set to sane values.
func New(port string, h http.Handler) Server {
	return Server{
		s: &http.Server{
			Addr:         net.JoinHostPort("0.0.0.0", port),
			Handler:      h,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout:  5 * time.Second,
		},
	}
}

// Addr returns the full address that the server is configured to listen on.
func (s Server) Addr() string {
	return s.s.Addr
}

// ListenAndServe listens for connections on interface 0.0.0.0 at the port
// provided to New. If the provided context is canceled, the server will attempt
// to gracefully shutdown. The returned error will only be non-nil if the server
// exits abnormally or fails to shutdown in time.
func (s Server) ListenAndServe(ctx context.Context) error {
	errC := make(chan error, 1)
	go func() { errC <- s.s.ListenAndServe() }()

	select {
	case <-ctx.Done():
		return s.shutdown()
	case err := <-errC:
		if err != http.ErrServerClosed {
			return err
		}
		return nil
	}
}

func (s Server) shutdown() error {
	// We have to create a new context here rather than passing one in from
	// ListenAndServe since that one was already Done'd by the time this method is
	// called.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.s.Shutdown(ctx)
}
