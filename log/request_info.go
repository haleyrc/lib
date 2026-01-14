package log

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

// ParseRequestInfo is a middleware that parses an incoming HTTP request and
// adds relevant fields to the request context. These fields are then added to
// any log lines emitted by one of the handler implementations in this package.
// This function also adds a unique identifier to the response for tracking
// across system boundaries.
func ParseRequestInfo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ri := &RequestInfo{
			ID:        getRequestID(r),
			IP:        getRealIP(r),
			Method:    r.Method,
			Path:      r.URL.Path,
			UserAgent: r.UserAgent(),
		}

		w.Header().Set("X-Request-ID", ri.ID)

		ctx := SetRequestInfo(r.Context(), ri)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequestInfo stores information about an HTTP request.
type RequestInfo struct {
	// A unique identifier for the request.
	ID string

	// The IP address of the upstream client.
	IP string

	// The HTTP method.
	Method string

	// The request path.
	Path string

	// The user agent string of the upstream client.
	UserAgent string
}

func getRequestInfo(ctx context.Context) *RequestInfo {
	tmp := ctx.Value(contextKey("request_info"))
	if tmp == nil {
		return nil
	}
	return tmp.(*RequestInfo)
}

// SetRequestInfo returns a copy of ctx with ri attached.
func SetRequestInfo(ctx context.Context, ri *RequestInfo) context.Context {
	return context.WithValue(ctx, contextKey("request_info"), ri)
}

// Depending on the path of the request through the network, the actual IP of
// the originating client can be located in a number of places. This function
// scans through the three most likely in order of specificity to return the
// best guess.
func getRealIP(r *http.Request) string {
	if ip := r.Header.Get("X-Real-Ip"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}

func getRequestID(r *http.Request) string {
	id := r.Header.Get("X-Request-ID")
	if id != "" {
		return id
	}
	return uuid.NewString()
}
