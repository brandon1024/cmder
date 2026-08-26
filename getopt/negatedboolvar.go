package getopt

import (
	"flag"
	"strconv"
)

// NegatedBoolVar is a boolean [flag.Value] for negating other flag values. A typical use case for NegatedBoolVar is to
// register flags which disable/unset other flags.
//
//	--show
//	--no-show
type NegatedBoolVar bool

// NegatedBool builds a [NegatedBoolVar] backed by b.
func NegatedBool(v *bool) *NegatedBoolVar {
	return (*NegatedBoolVar)(v)
}

// String returns the string representation of the (negated) boolean value.
func (b *NegatedBoolVar) String() string {
	return strconv.FormatBool(!bool(*b))
}

// Set updates the value of the flag. The given value must be a string parseable by [strconv.ParseBool]. The flag value
// is updated with the (negated) value s.
func (b *NegatedBoolVar) Set(s string) error {
	val, err := strconv.ParseBool(s)
	if err == nil {
		*b = NegatedBoolVar(!val)
	}

	return err
}

// IsBoolFlag marks the flag is not accepting args.
func (b *NegatedBoolVar) IsBoolFlag() bool {
	return true
}

// Get fulfills the [flag.Getter] interface, allowing typed access to the flag value. In this case, returns a bool.
func (b *NegatedBoolVar) Get() any {
	return !bool(*b)
}

// boolFlag is a [flag.Value] that also implements a method IsBoolFlag, used to determine if the flag accepts an
// argument or not.
type boolFlag interface {
	flag.Value
	IsBoolFlag() bool
}

// isBoolFlag checks if the given flag has a [flag.Value] which is a boolean flag.
func isBoolFlag(flg *flag.Flag) bool {
	bf, ok := flg.Value.(boolFlag)
	return ok && bf.IsBoolFlag()
}
