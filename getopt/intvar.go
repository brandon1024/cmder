package getopt

import (
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"
)

type (
	// UintType is a type constraint which describes all integer types accepted by [UintVar].
	UintType interface {
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
	}

	// UintVar is a generic [flag.Value] which can represent all native unsigned integer types. To initialize a UintVar,
	// see [Uint].
	//
	// Unlike the standard [flag.FlagSet] which currently only supports uint and uint64, UintVar supports all native
	// unsigned integer types (e.g. uint32).
	//
	//	var bin uint32
	//	fs.Var(getopt.Uint(&bin), "bytes", "the number of bytes to read")
	//
	// Types whose underlying type is an unsigned integer may also be used (for example, [fs.FileMode]).
	//
	//	var mode fs.FileMode
	//	fs.Var(getopt.Uint(&mode), "mode", "expected file mode permission bits (e.g. 0777)")
	//
	// In addition to simple numeric strings parsable by [strconv.ParseUint], UintVar can also parse numbers with
	// decimal or binary suffixes (e.g. 13G, 1.2Ki, 98.8M). The underlying type T must be suitably sized to represent
	// the result.
	//
	//	var bin uint32
	//	fs.Var(getopt.Uint(&bin), "bytes", "the number of bytes to read")
	//	fs.Parse([]string{"--bytes=3.2Mi"})               // bin = 3200000
	//
	// The following SI/IEC suffixes are recognized:
	//
	//	| SI (decimal)       | IEC (binary)       |
	//	| Value   | Suffix   | Value   | Suffix   |
	//	| ----------------------------------------|
	//	| 1000^1  | k kilo   | 1024^1  | Ki kibi  |
	//	| 1000^2  | M mega   | 1024^2  | Mi mebi  |
	//	| 1000^3  | G giga   | 1024^3  | Gi gibi  |
	//	| 1000^4  | T tera   | 1024^4  | Ti tebi  |
	//	| 1000^5  | P peta   | 1024^5  | Pi pebi  |
	//	| 1000^6  | E exa    | 1024^6  | Ei exbi  |
	//	| 1000^7  | Z zetta  | 1024^7  | Zi zebi  |
	//	| 1000^8  | Y yotta  | 1024^8  | Yi yobi  |
	//	| 1000^9  | R ronna  | 1024^9  | Ri robi  |
	//	| 1000^10 | Q quetta | 1024^10 | Qi quebi |
	UintVar[T UintType] struct {
		value *T
	}

	// IntType is a type constraint which describes all integer types accepted by [IntVar].
	IntType interface {
		~int | ~int8 | ~int16 | ~int32 | ~int64
	}

	// IntVar is a generic [flag.Value] which can represent all native signed integer types. To initialize an IntVar,
	// see [Int].
	//
	// Unlike the standard [flag.FlagSet] which currently only supports int and int64, IntVar supports all native
	// signed integer types (e.g. int32).
	//
	//	var bin int32
	//	fs.Var(getopt.Int(&bin), "bytes", "the number of bytes to read")
	//
	// Types whose underlying type is a signed integer may also be used (for example, [slog.Level]).
	//
	//	var lvl slog.Level
	//	fs.Var(getopt.Int(&lvl), "log-level", "log level for the application")
	//
	// In addition to simple numeric strings parsable by [strconv.ParseInt], IntVar can also parse numbers with
	// decimal or binary suffixes (e.g. 13G, 1.2Ki, 98.8M). The underlying type T must be suitably sized to represent
	// the result.
	//
	//	var bin int32
	//	fs.Var(getopt.Int(&bin), "bytes", "the number of bytes to read")
	//	fs.Parse([]string{"--bytes=3.2Mi"})               // bin = 3200000
	//
	// The following SI/IEC suffixes are recognized:
	//
	//	| SI (decimal)       | IEC (binary)       |
	//	| Value   | Suffix   | Value   | Suffix   |
	//	| ----------------------------------------|
	//	| 1000^1  | k kilo   | 1024^1  | Ki kibi  |
	//	| 1000^2  | M mega   | 1024^2  | Mi mebi  |
	//	| 1000^3  | G giga   | 1024^3  | Gi gibi  |
	//	| 1000^4  | T tera   | 1024^4  | Ti tebi  |
	//	| 1000^5  | P peta   | 1024^5  | Pi pebi  |
	//	| 1000^6  | E exa    | 1024^6  | Ei exbi  |
	//	| 1000^7  | Z zetta  | 1024^7  | Zi zebi  |
	//	| 1000^8  | Y yotta  | 1024^8  | Yi yobi  |
	//	| 1000^9  | R ronna  | 1024^9  | Ri robi  |
	//	| 1000^10 | Q quetta | 1024^10 | Qi quebi |
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
	if i.value == nil {
		panic("getopt: nil flag value")
	}

	v, err := i.parse(value)
	if err != nil {
		return err
	}

	*i.value = v

	return nil
}

// Get returns a pointer to the integer value within.
func (i *IntVar[T]) Get() any {
	return i.value
}

func (i *IntVar[T]) parse(value string) (T, error) {
	var zero T

	// first attempt, parse the value as a simple (decimal, hex) number
	v, err := strconv.ParseInt(value, 0, reflect.TypeOf(zero).Bits())
	if err == nil {
		return T(v), nil
	}

	// second attempt, parse the number as a float, optionally with an SI/IEC suffix (k, M, G, etc.)
	value, scale := cutSuffix(value)

	base, _, err := new(big.Float).Parse(value, 0)
	if err != nil {
		return zero, fmt.Errorf("getopt: cannot represent value '%s' as int", value)
	}

	number := new(big.Float).Mul(base, new(big.Float).SetInt(scale))
	v, acc := number.Int64()
	if acc != big.Exact {
		return zero, fmt.Errorf("getopt: cannot represent value '%s' (%s) as int", value, number.String())
	}

	if reflect.TypeOf(zero).OverflowInt(v) {
		return zero, fmt.Errorf("getopt: cannot represent value '%s' (%s) as int", value, number.String())
	}

	return T(v), nil
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
	if u.value == nil {
		panic("getopt: nil flag value")
	}

	v, err := u.parse(value)
	if err != nil {
		return err
	}

	*u.value = v

	return nil
}

// Get returns a pointer to the integer value within.
func (u *UintVar[T]) Get() any {
	return u.value
}

func (u *UintVar[T]) parse(value string) (T, error) {
	var zero T

	// first attempt, parse the value as a simple (decimal, hex) number
	v, err := strconv.ParseUint(value, 0, reflect.TypeOf(zero).Bits())
	if err == nil {
		return T(v), nil
	}

	// second attempt, parse the number as a float, optionally with an SI/IEC suffix (k, M, G, etc.)
	value, scale := cutSuffix(value)

	base, _, err := new(big.Float).Parse(value, 0)
	if err != nil {
		return zero, fmt.Errorf("getopt: cannot represent value '%s' as uint", value)
	}

	number := new(big.Float).Mul(base, new(big.Float).SetInt(scale))
	v, acc := number.Uint64()
	if acc != big.Exact {
		return zero, fmt.Errorf("getopt: cannot represent value '%s' (%s) as uint", value, number.String())
	}

	if reflect.TypeOf(zero).OverflowUint(v) {
		return zero, fmt.Errorf("getopt: cannot represent value '%s' (%s) as uint", value, number.String())
	}

	return T(v), nil
}

func cutSuffix(value string) (string, *big.Int) {
	si := map[string]int64{
		"k": 1, "M": 2, "G": 3, "T": 4, "P": 5, "E": 6, "Z": 7, "Y": 8, "R": 9, "Q": 10,
	}
	iec := map[string]int64{
		"Ki": 1, "Mi": 2, "Gi": 3, "Ti": 4, "Pi": 5, "Ei": 6, "Zi": 7, "Yi": 8, "Ri": 9, "Qi": 10,
	}

	if v, scale, ok := cutSuffixOf(value, si, 1000); ok {
		return v, scale
	}
	if v, scale, ok := cutSuffixOf(value, iec, 1024); ok {
		return v, scale
	}

	return value, big.NewInt(1)
}

func cutSuffixOf(value string, suffixes map[string]int64, base int64) (string, *big.Int, bool) {
	for suffix, exp := range suffixes {
		if v, ok := strings.CutSuffix(value, suffix); ok {
			return v, new(big.Int).Exp(big.NewInt(base), big.NewInt(exp), nil), true
		}
	}

	return value, nil, false
}
