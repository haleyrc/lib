// Package hash provides simple utilities for working with hashing algorithms.
package hash

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Check returns an error if the provided hash is not the hash of the provided
// guess or nil otherwise. The comparison is guaranteed to be constant time.
func Check(guess, hash string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(guess)); err != nil {
		return fmt.Errorf("check failed: %w", err)
	}
	return nil
}

// Generate returns a hashed version of the provided string. This function
// panics if there is an error, since there's not much that can be done and it
// simplifies the API significantly.
func Generate(s string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		// I'm panicking here because I don't think there's any way for this to
		// error that shouldn't immediately cause a page. I've never seen it happen
		// and I think it's just for interface satisfaction, so I feel safe here.
		// Plus, it's not a recoverable error. The user can't fix a broken hashing
		// algorithm.
		panic(err)
	}
	return string(hash)
}
