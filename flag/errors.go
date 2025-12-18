package flag

// ErrorHandling configures the error handling policy on a flag set, configuring the behaviour of [FlagSet.Parse].
type ErrorHandling int

const (
	// Return a descriptive error.
	ContinueOnError ErrorHandling = iota
	// Call os.Exit(2) or for -h/--help Exit(0).
	ExitOnError
	// Call panic with a descriptive error.
	PanicOnError
)
