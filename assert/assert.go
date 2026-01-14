// Package assert contains functions for making test assertions.
package assert

import (
	"net/http"
	"reflect"
	"slices"
	"strings"
)

// Result represents the result of an assertion nad is returned by all of the
// assertion functions in this package.
type Result struct {
	t      T
	failed bool
}

// ContentType validates that the value of the `Content-Type` header of the
// provided response matches the desired value.
func ContentType(t T, resp *http.Response, want string) Result {
	t.Helper()
	got := resp.Header.Get("Content-Type")
	if got != want {
		t.Errorf("Expected content type to be %s, but got %s.", want, got)
		return Result{t: t, failed: true}
	}
	return Result{t: t, failed: false}
}

// DeepEqual validates that two values are "deeply equal" according to the same
// rules as [reflect.DeepEqual].
//
// This method is best when comparing e.g. two instances of the same type where
// you are concerned about the equality of their values and not whether they are
// the same object in memory. For instance, given the following objects:
//
//	type Person struct {
//		Name string
//	}
//
//	superman1 := Person{Name: "Superman"}
//	superman2 := Person{Name: "Superman"}
//
// the following assertions are true:
//
//	assert.Equal(t, "supermen", superman1, superman2)
//	assert.DeepEqual(t, "supermen", superman1, superman2)
//
// since we are passing the objects by value and they have the same values for
// their fields. To contrast, of the following assertions:
//
//	assert.Equal(t, "superman", &superman1, &superman2)
//	assert.DeepEqual(t, "superman", &superman1, &superman2)
//
// only the DeepEqual assertions succeeds since the call to Equal compares the
// values of the pointers, which are different for different instances.
func DeepEqual(t T, label string, want, got any) Result {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Expected %s to be equal, but they weren't.", label)
		return Result{t: t, failed: true}
	}
	return Result{t: t, failed: false}
}

// Equal validates that two values are the same.
//
// This method is best when comparing "simple" types e.g. int, string, etc. and
// CAN be used to compare two objects passed by value or two objects passed by
// reference where you care whether the references are the same object. For
// instance, given the following objects:
//
//	type Person struct {
//		Name string
//	}
//
//	superman1 := Person{Name: "Superman"}
//	superman2 := Person{Name: "Superman"}
//
// the following assertions are true:
//
//	assert.Equal(t, "supermen", superman1, superman2)
//	assert.DeepEqual(t, "supermen", superman1, superman2)
//
// since we are passing the objects by value and they have the same values for
// their fields. To contrast, of the following assertions:
//
//	assert.Equal(t, "superman", &superman1, &superman2)
//	assert.DeepEqual(t, "superman", &superman1, &superman2)
//
// only the DeepEqual assertions succeeds since the call to Equal compares the
// values of the pointers, which are different for different instances.
func Equal[C comparable](t T, label string, want, got C) Result {
	t.Helper()
	if got != want {
		t.Errorf("Expected %s to be %v, but got %v.", label, want, got)
		return Result{t: t, failed: true}
	}
	return Result{t: t, failed: false}
}

// Error validates that the provided error is not nil and contains the desired
// string. Note that the comparison is equivalent to [strings.Contains], so the
// following assertion succeeds:
//
//	assert.Error(t, errors.New("oops: invalid"), "invalid")
func Error(t T, err error, want string) Result {
	t.Helper()
	if err == nil {
		t.Errorf("Expected error to not be nil, but it was.")
		return Result{t: t, failed: true}
	}

	got := err.Error()
	if !strings.Contains(got, want) {
		t.Errorf("Expected error to contain %q, but got %q.", want, got)
		return Result{t: t, failed: true}
	}

	return Result{t: t, failed: false}
}

// False validates that the provided value is false.
func False(t T, label string, got bool) Result {
	t.Helper()
	return Equal(t, label, false, got)
}

// NotBlank validates that the provided string is not the blank string. Leading
// and trailing spaces are removed from got before validation.
func NotBlank(t T, label string, got string) Result {
	t.Helper()
	got = strings.TrimSpace(got)
	if got == "" {
		t.Errorf("Expected %s to not be blank, but it was.", label)
		return Result{t: t, failed: true}
	}
	return Result{t: t, failed: false}
}

// OK validates that the provided err is nil.
func OK(t T, err error) Result {
	t.Helper()
	if err != nil {
		t.Errorf("Unexpected error: %v.", err)
		return Result{t: t, failed: true}
	}
	return Result{t: t, failed: false}
}

// ShouldPanic validates that calling f results in a panic. This can be useful
// when testing methods that panic on error rather than returning a value (and
// these should be restricted to the types of things called from the main
// package).
func ShouldPanic(t T, f func()) (result Result) {
	t.Helper()
	result = Result{t: t}
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("Expected function to panic, but it didn't.")
			result.failed = true
			return
		}
	}()
	f()
	return
}

// SliceEqual validates that two slices are the same. This function does not
// modify the provided slices in any way, so you may need to sort both inputs
// prior to comparison.
func SliceEqual[S ~[]E, E comparable](t T, label string, want, got S) Result {
	t.Helper()

	if !slices.Equal(got, want) {
		t.Errorf("Expected %s to be %v, but got %v.", label, want, got)
		return Result{t: t, failed: true}
	}

	return Result{t: t, failed: false}
}

// StatusCode validates that the status code of the provided response matches
// the desired value.
func StatusCode(t T, want int, resp *http.Response) Result {
	t.Helper()
	got := resp.StatusCode
	if got != want {
		t.Errorf("Expected status code to be %d, but got %d.", want, got)
		return Result{t: t, failed: true}
	}
	return Result{t: t, failed: false}
}

// True validates that the provided value is true.
func True(t T, label string, got bool) Result {
	t.Helper()
	return Equal(t, label, true, got)
}

// Fatal causes the test suite to immediately fail if the current result
// corresponds to a failed assertion. You can chain this off of any of the
// assertion functions and is most often useful for exiting due to an unexpected
// error, e.g.:
//
//	assert.OK(t, err).Fatal()
func (r Result) Fatal() {
	r.t.Helper()
	if r.failed {
		r.t.FailNow()
	}
}

// OK returns true if the current result corresponds to a failed assertion or
// false otherwise.
func (r Result) OK() bool {
	r.t.Helper()
	return !r.failed
}

// T wraps the basic testing methods.
//
// All assertions in this package take a T as a first argument. In real
// use-cases, this will almost always be an instance of [testing.T] e.g.:
//
//	func TestSomeThing(t *testing.T) {
//		assert.True(t, "identity", true)
//	}
type T interface {
	Helper()
	Errorf(format string, args ...any)
	FailNow()
	Log(args ...any)
}
