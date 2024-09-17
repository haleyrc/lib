// Package web provides utilities for constructing HTTP responses.
package web

import (
	"encoding/json"
	"net/http"
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

// StatusCode sets the HTTP status code.
func StatusCode(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
}
