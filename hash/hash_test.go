package hash_test

import (
	"testing"

	"github.com/haleyrc/lib/assert"
	"github.com/haleyrc/lib/hash"
)

func TestGenerate(t *testing.T) {
	h := hash.New("hello")
	assert.OK(t, h.Compare("hello"))
	assert.Error(t, h.Compare("goodbye"), "hash mismatch")
}
