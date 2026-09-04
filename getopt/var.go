package getopt

import (
	"fmt"
	"reflect"
	"strconv"
)

// IntType is a type constraint which describes all integer types accepted by [IntVar].
type IntType interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// UintType is a type constraint which describes all integer types accepted by [UintVar].
type UintType interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type (
	// UintVar is a generic [flag.Value] which can represent all native unsigned integer types. There are a few cases
	// where UintVar can be useful:
	//
	// 	- A UintVar can be provided to Var() of [flag.FlagSet].
	// 	- With the help of generics, UintVar supports all unsigned integer types. The standard [flag.FlagSet] currently
	//    only supports uint and uint64.
	// 	- User-defined types whose underlying type is an unsigned integer may also be used (for example, [net.Flags]).
	//
	// To initialize a UintVar, see [Uint].
	UintVar[T UintType] struct {
		value *T
	}

	// IntVar is a generic [flag.Value] which can represent all native signed integer types. There are a few cases
	// where IntVar can be useful:
	//
	// 	- An IntVar can be provided to Var() of [flag.FlagSet].
	// 	- With the help of generics, IntVar supports all signed integer types. The standard [flag.FlagSet] currently
	//    only supports int and int64.
	// 	- User-defined types whose underlying type is a signed integer may also be used (for example, [slog.Level]).
	//
	// To initialize an IntVar, see [Int].
	IntVar[T IntType] struct {
		value *T
	}
)

// Int initializes an [IntVar] wrapping value.
func Int[T IntType](value *T) *IntVar[T] {
	return &IntVar[T]{
		value: value,
	}
}

// String returns the string decimal representation of the signed integer value within.
func (i *IntVar[T]) String() string {
	var zero T

	if i == nil || i.value == nil {
		return fmt.Sprintf("%d", zero)
	}

	return fmt.Sprintf("%d", *i.value)
}

// Set updates the integer value of i. Set accepts any string which can be parsed by [strconv.ParseInt] within the
// bounds of the generic type T.
func (i *IntVar[T]) Set(value string) error {
	if i == nil || i.value == nil {
		panic("getopt: nil flag value")
	}

	var zero T

	v, err := strconv.ParseInt(value, 0, reflect.TypeOf(zero).Bits())
	if err != nil {
		return err
	}

	*i.value = T(v)

	return nil
}

// Get returns a pointer to the integer value within.
func (i *IntVar[T]) Get() any {
	return i.value
}

// Uint initializes a [UintVar] wrapping value.
func Uint[T UintType](value *T) *UintVar[T] {
	return &UintVar[T]{
		value: value,
	}
}

// String returns the string decimal representation of the unsigned integer value within.
func (u *UintVar[T]) String() string {
	var zero T

	if u == nil || u.value == nil {
		return fmt.Sprintf("%d", zero)
	}

	return fmt.Sprintf("%d", *u.value)
}

// Set updates the integer value of u. Set accepts any string which can be parsed by [strconv.ParseUint] within the
// bounds of the generic type T.
func (u *UintVar[T]) Set(value string) error {
	if u == nil || u.value == nil {
		panic("getopt: nil flag value")
	}

	var zero T

	v, err := strconv.ParseUint(value, 0, reflect.TypeOf(zero).Bits())
	if err != nil {
		return err
	}

	*u.value = T(v)

	return nil
}

// Get returns a pointer to the integer value within.
func (u *UintVar[T]) Get() any {
	return u.value
}
