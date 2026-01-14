// Package hash provides functions for computing and comparing one-way hashes.
package hash

import (
	"fmt"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

// Hash represents a one-way hash of a string value.
type Hash struct {
	value []byte
}

// New returns a one-way hash of s.
func New(s string) Hash {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		// This can only fail in really arcane ways that we have no suitable way of
		// recovering from, so we just panic here rather than trying to handle the
		// error.
		panic(err)
	}

	hash := Hash{
		value: hashBytes,
	}

	return hash
}

// Compare compares o with the hashed value in h and returns an error if the
// values don't match.
func (h Hash) Compare(o string) error {
	if err := bcrypt.CompareHashAndPassword(h.value, []byte(o)); err != nil {
		return fmt.Errorf("hash: compare: hash mismatch: %w", err)
	}
	return nil
}

func (h Hash) LogValue() slog.Value {
	return slog.StringValue("SECRET")
}
