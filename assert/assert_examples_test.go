package assert_test

import (
	"errors"
	"net/http"

	"github.com/haleyrc/lib/assert"
)

func ExampleContentType() {
	resp := new(http.Response)

	header := http.Header{}
	header.Set("Content-Type", "application/json")
	resp.Header = header

	assert.ContentType(t, resp, "application/json")
	assert.ContentType(t, resp, "application/xml")

	// Output: Expected content type to be application/xml, but got application/json.
}

func ExampleDeepEqual() {
	type Composer struct {
		Name string
	}

	bach1 := Composer{Name: "J.S. Bach"}
	bach2 := Composer{Name: "J.S. Bach"}
	shostakovich := Composer{Name: "D. Shostakovich"}

	// Objects passed by value are considered identical if the values of their
	// fields are identical.
	assert.DeepEqual(t, "composers", bach1, bach2)
	assert.DeepEqual(t, "composers", bach1, shostakovich)

	// Objects passed by reference are ALSO considered identical if the values of
	// their fields are identical, regardless of which object is being referenced.
	assert.DeepEqual(t, "composers", &bach1, &bach1)
	assert.DeepEqual(t, "composers", &bach1, &bach2)
	assert.DeepEqual(t, "composers", &bach1, &shostakovich)

	// Output: Expected composers to be equal, but they weren't.
	// Expected composers to be equal, but they weren't.
}

func ExampleEqual_complexTypes() {
	type Robot struct {
		Name string
	}

	bender1 := Robot{Name: "Bender"}
	bender2 := Robot{Name: "Bender"}
	flexo := Robot{Name: "Flexo"}

	// Objects passed by value are considered equal if the values of their fields
	// are equal.
	assert.Equal(t, "robot", bender1, bender2)
	assert.Equal(t, "robot", bender1, flexo)

	// Objects passed by reference are considered equal iff they refer to the same
	// object.
	assert.Equal(t, "robot", &bender1, &bender1)
	assert.Equal(t, "robot", &bender1, &bender2)

	// Output: Expected robot to be {Bender}, but got {Flexo}.
	// Expected robot to be &{Bender}, but got &{Bender}.
}

func ExampleEqual_simpleTypes() {
	assert.Equal(t, "int", 42, 42)
	assert.Equal(t, "int", 42, 13)

	assert.Equal(t, "string", "Hello, World", "Hello, World")
	assert.Equal(t, "string", "Hello, World", "Goodbye, World")

	// Output: Expected int to be 42, but got 13.
	// Expected string to be Hello, World, but got Goodbye, World.
}

func ExampleError() {
	err := errors.New("oops: invalid syntax")
	assert.Error(t, err, "oops")
	assert.Error(t, err, "invalid syntax")
	assert.Error(t, err, "oops: invalid syntax")
	assert.Error(t, err, "invalid sintacks")

	// Output: Expected error to contain "invalid sintacks", but got "oops: invalid syntax".
}

func ExampleFalse() {
	assert.False(t, "true", true)
	assert.False(t, "false", false)

	// Output: Expected true to be false, but got true.
}

func ExampleNotBlank() {
	assert.NotBlank(t, "the blank string", "")
	assert.NotBlank(t, "only spaces", "    ")
	assert.NotBlank(t, "leading spaces", "   Hello")
	assert.NotBlank(t, "trailing spaces", "World   ")
	assert.NotBlank(t, "no spaces", "Hello")

	// Output: Expected the blank string to not be blank, but it was.
	// Expected only spaces to not be blank, but it was.
}

func ExampleOK() {
	assert.OK(t, nil)
	assert.OK(t, errors.New("oops"))

	// Output: Unexpected error: oops.
}

func ExampleShouldPanic() {
	assert.ShouldPanic(t, func() {})
	assert.ShouldPanic(t, func() {
		panic("oops")
	})

	// Output: Expected function to panic, but it didn't.
}

func ExampleSliceEqual() {
	control := []int{1, 2, 3}
	reversed := []int{3, 2, 1}
	longer := []int{1, 2, 3, 4}
	wildcard := []int{42, 007, 2716057}

	assert.SliceEqual(t, "the identity", control, control)
	assert.SliceEqual(t, "out of order elements", control, reversed)
	assert.SliceEqual(t, "different lengths", control, longer)
	assert.SliceEqual(t, "different elements", control, wildcard)

	assert.SliceEqual(t, "struct elements",
		[]struct{ Value int }{{1}, {2}, {3}},
		[]struct{ Value int }{{1}, {2}, {3}},
	)
	assert.SliceEqual(t, "struct elements",
		[]struct{ Value int }{{1}, {2}, {3}},
		[]struct{ Value int }{{3}, {1}, {2}},
	)

	// Output: Expected out of order elements to be [1 2 3], but got [3 2 1].
	// Expected different lengths to be [1 2 3], but got [1 2 3 4].
	// Expected different elements to be [1 2 3], but got [42 7 2716057].
	// Expected struct elements to be [{1} {2} {3}], but got [{3} {1} {2}].
}

func ExampleStatusCode() {
	resp := new(http.Response)
	resp.StatusCode = 200

	assert.StatusCode(t, http.StatusOK, resp)
	assert.StatusCode(t, http.StatusTeapot, resp)

	// Output: Expected status code to be 418, but got 200.
}

func ExampleTrue() {
	assert.True(t, "true", true)
	assert.True(t, "false", false)

	// Output: Expected false to be true, but got false.
}
