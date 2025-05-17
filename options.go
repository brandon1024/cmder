package cmder

// Options used to configure behaviour of [Execute].
type ExecuteOptions struct {
	args []string
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
