package flag

// Getter extends [Value] but additionally supports retrieval of the typed flag value. All flag types implement this
// interface, with the exception of [FlagSet.BoolFunc] and [FlagSet.Func].
type Getter interface {
	Value

	// Get yields the (typed) flag value.
	Get() any
}
