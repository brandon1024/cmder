// Package tutil offers simple utilities for use in tests.
package tutil

import (
	"errors"
	"fmt"
	"slices"
	"testing"
)

// Result represents the result of an assertion.
type Result struct {
	// res is false if the assertion failed.
	res bool
	// msg is a descriptive message for the failed assertion.
	msg string
}

// Assert fails the test if assertion res failed.
func Assert(t *testing.T, res Result) {
	if !res.res {
		t.Fatalf("expectation failed: %s", res.msg)
	}
}

// Eq asserts that the given values are equal (==).
func Eq[T comparable](expected, actual T) Result {
	return Result{expected == actual, fmt.Sprintf("values not equal: expected %v but was %v", expected, actual)}
}

// NilErr asserts that the given error is nil.
func NilErr(err error) Result {
	return Result{err == nil, fmt.Sprintf("unexpected error: %v", err)}
}

// IsErr asserts that the given error matches target (with [errors.Is]).
func IsErr(err, target error) Result {
	return Result{errors.Is(err, target), fmt.Sprintf("unexpected error: %v", err)}
}

// Match asserts that the given slices have the same values.
func Match[S ~[]E, E comparable](expected, actual S) Result {
	return Result{slices.Equal(expected, actual), fmt.Sprintf("slices not equal: expected %v but was %v", expected, actual)}
}
