// Package httputil provides basic utilities for working with and inspecting
// HTTP requests and responses.
package httputil

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
)

// RequestID is the name of the header expected to contain a unique identifier
// for incoming requests and their matching outgoing responses.
const RequestID = "X-Request-ID"

// // DumpRequest prints r to stderr including headers and the request body. The
// request is reset after printing so it's safe to reuse. If the request can't
// be dumped, an error is emitted to stderr instead.
func DumpRequest(r *http.Request) {
	bytes, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dump request failed: %v\n", err)
		return
	}
	fmt.Fprintln(os.Stderr, string(bytes))
}

// DumpResponse prints resp to stderr including the headers and body. The
// response is reset after printing so it's safe to reuse. If the response can't
// be dumped, an error is emitted to stderr instead.
func DumpResponse(resp *http.Response) {
	bytes, err := httputil.DumpResponse(resp, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dump response failed: %v\n", err)
		return
	}
	fmt.Fprintln(os.Stderr, string(bytes))
}

// RealIP returns the best guess at the originating IP address for an incoming
// request.
func RealIP(r *http.Request) string {
	if ip := r.Header.Get("X-Real-Ip"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}
