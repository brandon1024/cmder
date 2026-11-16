package cmder

import (
	"flag"
	"io"
)

// ExecuteOptions configure the behavior of [Execute].
type ExecuteOptions struct {
	args          []string
	nativeFlags   bool
	relaxedFlags  bool
	bindEnv       bool
	bindEnvMapper func([]Command, Command, *flag.Flag) (string, bool)
	interspersed  bool

	usageTemplate      string
	helpTemplate       string
	secondaryTemplates map[string]string
	outputWriter       io.Writer
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

// WithRelaxedFlagParsing instructs [Execute] to relax flag parsing. Partial name matches for long flags are permitted,
// granted the provided flag name is not ambiguous.
//
// This option is ignored if [WithNativeFlags] is enabled.
func WithRelaxedFlagParsing() ExecuteOption {
	return func(ops *ExecuteOptions) {
		ops.relaxedFlags = true
	}
}

// WithEnvironmentBinding configures [Execute] to set flag values from the enclosing environment. Environment variables
// are mapped to flags as follows:
//
//	COMMAND_FLAGNAME
//	COMMAND_SUBCOMMAND_FLAGNAME
//	COMMAND_SUBCOMMAND_SUBCOMMAND_FLAGNAME
//
// Command and flag names are made uppercase. Special characters are removed. Flags explicitly set at the command line
// take precedence over environment variables.
//
//	git log --format=oneline   ->   GIT_LOG_FORMAT=oneline
//	git log --no-abbrev-commit ->   GIT_LOG_NOABBREVCOMMIT=true
//
// For greater control, see [WithCustomEnvironmentBinding].
func WithEnvironmentBinding() ExecuteOption {
	return func(ops *ExecuteOptions) {
		ops.bindEnv = true
		ops.bindEnvMapper = defaultEnvironmentMapper
	}
}

// WithCustomEnvironmentBinding is like [WithEnvironmentBinding] but with a user-defined mapping function for mapping
// command flags to environment variable names.
//
// The mapping function accepts a slice of parent commands (nil if the command has no parents), the command itself, and
// a flag. The mapper must return a valid environment variable name. If the mapping function returns false, environment
// variable binding is disabled for the flag.
//
// See also [WithEnvironmentBinding].
func WithCustomEnvironmentBinding(mapper func([]Command, Command, *flag.Flag) (string, bool)) ExecuteOption {
	return func(ops *ExecuteOptions) {
		ops.bindEnv = true
		ops.bindEnvMapper = mapper
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

// WithHelpTemplate is used to provide an alternate template for rendering command help text. The template is
// rendered by the standard [text/template] package. This is particularly useful for applications which prefer to format
// command help text differently than the cmder defaults.
//
// By default, the [DefaultHelpTemplate] template is used.
//
// See also [WithUsageTemplate] and [WithOutputWriter].
func WithHelpTemplate(tmpl string) ExecuteOption {
	return func(ops *ExecuteOptions) {
		ops.helpTemplate = tmpl
	}
}

// WithUsageTemplate is used to provide an alternate template for rendering command usage text. The template is
// rendered by the standard [text/template] package. This is particularly useful for applications which prefer to format
// command usage information differently than the cmder defaults.
//
// By default, the [DefaultUsageTemplate] template is used.
//
// See also [WithHelpTemplate] and [WithOutputWriter].
func WithUsageTemplate(tmpl string) ExecuteOption {
	return func(ops *ExecuteOptions) {
		ops.usageTemplate = tmpl
	}
}

// WithNamedTemplate registers an additional template with the given name. Additional templates are made available to
// the base usage/help template, which is rendered with the standard [text/template] package.
//
// A common use case is to selectively define or override sections of the rendered output by using {{define ...}},
// {{block ...}} and {{template ...}} statements.
//
// See also [WithHelpTemplate] and [WithUsageTemplate] to override base templates.
func WithNamedTemplate(name, tmpl string) ExecuteOption {
	return func(ops *ExecuteOptions) {
		ops.secondaryTemplates[name] = tmpl
	}
}

// WithOutputWriter is used to provide an alternate [io.Writer] to write rendered command usage/help text. By default,
// [os.Stdout] is used.
//
// See also [WithHelpTemplate] and [WithUsageTemplate].
func WithOutputWriter(output io.Writer) ExecuteOption {
	return func(ops *ExecuteOptions) {
		ops.outputWriter = output
	}
}
