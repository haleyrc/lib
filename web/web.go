// Package web provides utilities for constructing HTTP responses.
package web

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

const (
	ContentTypeApplicationJSON = "application/json"
)

// ContentType sets the Content-Type header.
func ContentType(w http.ResponseWriter, ct string) {
	Header(w, "Content-Type", ct)
}

// Header sets a header key to a single value.
func Header(w http.ResponseWriter, key, value string) {
	w.Header().Set(key, value)
}

// JSON writes the JSON representation of the provided data to a response.
func JSON(w http.ResponseWriter, body any) {
	bytes, err := json.MarshalIndent(body, "", "  ")
	if err != nil {
		// I don't love this, but JSON marshaling errors should be caught in tests
		// and by not returning an error, we're making it easier to find issues that
		// might go unnoticed if tests are spotty and a returned error isn't
		// checked.
		panic(err)
	}
	w.Write(bytes)
}

// NewInt converts a string to an int64 and is primarily used for parsing
// incoming form values.
func NewInt(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

// NewString sanitizes its input and returns the result. This is primarily used
// for parsing incoming form values.
func NewString(s string) string {
	s = strings.TrimSpace(s)
	return s
}

// StatusCode sets the HTTP status code.
func StatusCode(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
}
