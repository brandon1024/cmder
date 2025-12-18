package flag

// Value describes the actual value of a [Flag].
type Value interface {
	// String returns the current value of the flag, represented as a string.
	String() string

	// Set updates the value of this flag.
	Set(string) error
}
