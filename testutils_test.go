package cmder

import (
	"fmt"
	"slices"
	"testing"
)

// result represents the result of an assertion. res is false if the assertion failed. msg is a descriptive message for
// the failed assertion.
type result struct {
	res bool
	msg string
}

// assert fails the test if assertion res failed.
func assert(t *testing.T, res result) {
	if !res.res {
		t.Fatalf("expectation failed: %s", res.msg)
	}
}

// eq asserts that the given values are equal (==).
func eq[T comparable](expected, actual T) result {
	return result{expected == actual, fmt.Sprintf("values not equal: expected %v but was %v", expected, actual)}
}

// nilerr asserts that the given error is nil.
func nilerr(err error) result {
	return result{err == nil, fmt.Sprintf("unexpected error: %v", err)}
}

// match asserts that the given slices have the same values.
func match[S ~[]E, E comparable](expected, actual S) result {
	return result{slices.Equal(expected, actual), fmt.Sprintf("slices not equal: expected %v but was %v", expected, actual)}
}
