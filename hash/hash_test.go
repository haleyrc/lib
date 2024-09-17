package hash_test

import (
	"testing"

	"github.com/haleyrc/lib/hash"
)

func TestGenerate(t *testing.T) {
	original := "mystring"
	hashed := hash.Generate(original)

	if err := hash.Check(original, hashed); err != nil {
		t.Errorf("Expected %q to be the hash of %q, but got error %v.", hashed, original, err)
	}
}
