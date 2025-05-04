package cmder

import "flag"

// Options used to configure behaviour of [Execute].
type ExecuteOptions struct {
	args  []string
	flags *flag.FlagSet
}

// A single option passed to [Execute].
type ExecuteOption func(*ExecuteOptions)

// WithArgs configures [Execute] to run with the arguments given. By default, [Execute] will execute with arguments from
// [os.Args].
func WithArgs(args []string) ExecuteOption {
	return func(ops *ExecuteOptions) {
		ops.args = args
	}
}

// WithFlags configure [Execute] to run with a specific [flag.FlagSet]. This set of flags will only be used when
// executing the top-level command passed to [Execute], and will not be passed down to subcommands.
func WithFlags(flags *flag.FlagSet) ExecuteOption {
	return func(ops *ExecuteOptions) {
		ops.flags = flags
	}
}
