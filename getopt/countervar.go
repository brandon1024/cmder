package getopt

import (
	"fmt"
	"strconv"
)

// CounterType describes all numeric types supported by the [CounterVar] type.
type CounterType interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// CounterVar is a boolean [flag.Value] which increments a signed or unsigned integer of type T every time it appears at
// the command line.
//
// This type of flag is often used for describing log verbosity, where '-vvv' would be interpreted as setting the log
// level to 3.
//
// To initialize a CounterVar, see [Counter].
type CounterVar[T CounterType] struct {
	value *T
}

// Counter Initializes a [CounterVar] with an initial value.
func Counter[T CounterType](value *T) *CounterVar[T] {
	return &CounterVar[T]{
		value: value,
	}
}

// String returns the value of the counter as a string.
func (c *CounterVar[T]) String() string {
	var zero T

	if c == nil || c.value == nil {
		return fmt.Sprintf("%d", zero)
	}

	return fmt.Sprintf("%d", *c.value)
}

// Set accepts a boolean value. If true, the counter is incremented.
func (c *CounterVar[T]) Set(value string) error {
	if c == nil || c.value == nil {
		panic("getopt: nil flag value")
	}

	v, err := strconv.ParseBool(value)
	if err != nil {
		return err
	}

	if v {
		*c.value++
	}

	return nil
}

// Get returns the counter value as an integer.
func (c *CounterVar[T]) Get() any {
	return c.value
}

// IsBoolFlag marks the flag as being a boolean.
func (c *CounterVar[T]) IsBoolFlag() bool {
	return true
}
