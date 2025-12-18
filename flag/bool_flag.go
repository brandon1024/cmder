package flag

// boolFlag is a [Value] that also implements a method IsBoolFlag, used to determine if the flag accepts an argument or
// not.
type boolFlag interface {
	Value
	IsBoolFlag() bool
}
