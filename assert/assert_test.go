package assert_test

import (
	"fmt"
	"os"
)

// N.B.: These definitions need to exist in a separate file from the testable
// examples to prevent the documentation from including them in every example
// block.

var t mockT

type mockT struct{}

func (mockT) Errorf(format string, args ...any) {
	fmt.Fprintf(os.Stdout, format, args...)
	fmt.Fprintln(os.Stdout)
}

func (mockT) FailNow() {}

func (mockT) Helper() {}

func (mockT) Log(args ...any) {
	fmt.Fprintln(os.Stdout, args...)
}
