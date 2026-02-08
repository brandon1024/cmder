package cmder

// ExecuteOptions configure the behaviour of [Execute].
type ExecuteOptions struct {
	args          []string
	nativeFlags   bool
	bindEnv       bool
	bindEnvPrefix string
	interspersed  bool
}

// ExecuteOption is a single option passed to [Execute].
type ExecuteOption func(*ExecuteOptions)

// WithArgs configures [Execute] to run with the arguments given. By default, [Execute] will execute with arguments from
// [os.Args].
func WithArgs(args []string) ExecuteOption {
	return func(ops *ExecuteOptions) {
		ops.args = args
	}
}

// WithNativeFlags configures [Execute] to parse flags using the standard [flag] package instead of the default
// [getopt] package.
func WithNativeFlags() ExecuteOption {
	return func(ops *ExecuteOptions) {
		ops.nativeFlags = true
	}
}

// WithEnvironmentBinding configures [Execute] to set flag values from the enclosing environment. Environment variables
// are mapped to flags as follows:
//
//	COMMAND_FLAGNAME
//	COMMAND_SUBCOMMAND_FLAGNAME
//	COMMAND_SUBCOMMAND_SUBCOMMAND_FLAGNAME
//
// Command and flag names are all uppercased. Special characters are removed. Flags explicitly set at the command line
// take precedence over environment variables.
//
//	git log --format=oneline   ->   GIT_LOG_FORMAT=oneline
//	git log --no-abbrev-commit ->   GIT_LOG_NOABBREVCOMMIT=true
func WithEnvironmentBinding() ExecuteOption {
	return func(ops *ExecuteOptions) {
		ops.bindEnv = true
		ops.bindEnvPrefix = ""
	}
}

// WithPrefixedEnvironmentBinding is like [WithEnvironmentBinding] but with a variable name prefix.
//
//	<PREFIX>COMMAND_FLAGNAME
//	<PREFIX>COMMAND_SUBCOMMAND_FLAGNAME
//	<PREFIX>COMMAND_SUBCOMMAND_SUBCOMMAND_FLAGNAME
func WithPrefixedEnvironmentBinding(prefix string) ExecuteOption {
	return func(ops *ExecuteOptions) {
		ops.bindEnv = true
		ops.bindEnvPrefix = prefix
	}
}

// WithInterspersedArgs enables interspersed args parsing, allowing command-line arguments and flags to be mixed. When
// interspersed arg parsing is enabled, the following is permitted:
//
//	git log origin/main -p
//
// When interspersed arg parsing is disabled, flags must always come before args:
//
//	git log -p origin/main
func WithInterspersedArgs() ExecuteOption {
	return func(ops *ExecuteOptions) {
		ops.interspersed = true
	}
}
