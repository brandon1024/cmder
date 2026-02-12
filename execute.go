package cmder

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/brandon1024/cmder/getopt"
)

// ErrIllegalCommandConfiguration is an error returned when a [Command] provided to [Execute] is illegal.
var ErrIllegalCommandConfiguration = errors.New("cmder: illegal command configuration")

// ErrEnvironmentBindFailure is an error returned when [Execute] failed to update a flag value from environment
// variables (see [WithEnvironmentBinding]).
var ErrEnvironmentBindFailure = errors.New("cmder: failed to update flag from environment variable")

// Execute runs a [Command].
//
// # Execution Lifecycle
//
// When executing a command, Execute will call the [Runnable] Run() routine of your command. If the command also
// implements [RunnableLifecycle], the [RunnableLifecycle] Initialize() and Destroy() routines will be invoked before
// and after calling Run().
//
// If the command implements [RootCommand] and a subcommand is invoked, Execute will invoke the [RunnableLifecycle]
// routines of parent and child commands:
//
//  1. Root  [RunnableLifecycle] Initialize()
//  2. Child [RunnableLifecycle] Initialize()
//  3. Child [Runnable] Run()
//  4. Child [RunnableLifecycle] Destroy()
//  5. Root  [RunnableLifecycle] Destroy()
//
// If a command implements [RootCommand] but the first argument passed to the command doesn't match a recognized child
// command Name(), the Run() routine will be executed.
//
// # Error Handling
//
// Whenever a lifecycle routine (Initialize(), Run(), Destroy()) returns a non-nil error, execution is aborted
// immediately and the error is returned at once. For example, returning an error from Run() will prevent execution of
// Destroy() of the current command and any parents.
//
// Execute may return [ErrIllegalCommandConfiguration] if a command is misconfigured.
//
// # Command Contexts
//
// A [context.Context] derived from ctx is passed to all lifecycle routines. The context is cancelled when Execute
// returns. Commands should use this context to manage their resources correctly.
//
// # Execution Options
//
// Execute accepts one or more [ExecuteOption] options. You can provide these options to tweak the behaviour of Execute.
//
// # Flag Initialization
//
// If the command also implements [FlagInitializer], InitializeFlags() will be invoked to register additional
// command-line flags. Each command/subcommand is given a unique [flag.FlagSet]. Help flags ('-h', '--help') are
// configured automatically and must not be set by the application.
//
// Execute parses getopt-style (GNU/POSIX) command-line arguments with the help of package [getopt]. To use the standard
// [flag] syntax instead, see [WithNativeFlags]. Flags and arguments cannot be interspersed by default. You can change
// this behaviour with [WithInterspersedArgs].
//
// To bind environment variables to flags, see [WithEnvironmentBinding].
//
// # Usage and Help Texts
//
// Whenever the user provides the '-h' or '--help' flag at the command line, [Execute] will display command usage and
// exit. The format of the help text can be adjusted with [WithUsageTemplate]. By default, usage information will
// be written to stderr, but this can be adjusted by setting [WithUsageOutput].
//
// If a command's [Run] routine returns [ErrShowUsage] (or an error wrapping [ErrShowUsage]), [Execute] will render
// help text and exit with status 2.
func Execute(ctx context.Context, cmd Command, op ...ExecuteOption) error {
	// do some checks
	if cmd == nil {
		return errors.Join(ErrIllegalCommandConfiguration, errors.New("cmder: command cannot be nil"))
	}

	// prepare executor options
	ops := &ExecuteOptions{
		args:          os.Args[1:],
		usageTemplate: CobraUsageTemplate,
		usageWriter:   os.Stderr,
	}
	for _, f := range op {
		f(ops)
	}

	// build a stack of command invocations
	stack, err := buildCallStack(cmd, ops)
	if err != nil {
		return err
	}

	// if help was requested, display and exit
	if cmd, ok := helpRequested(stack); ok {
		return usage(*cmd, ops)
	}

	return execute(ctx, stack, ops)
}

// execute traverses the command stack recursively executing the lifecycle routines at each level.
func execute(ctx context.Context, stack []command, ops *ExecuteOptions) error {
	if len(stack) == 0 {
		return nil
	}

	// setup context
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var (
		this = stack[0]
		err  error
	)

	// run init (if applicable)
	if err := this.onInit(ctx, ops); err != nil {
		return err
	}

	// if this is a leaf, run, otherwise recurse
	if len(stack) == 1 {
		err = this.run(ctx, ops)
	} else {
		err = execute(ctx, stack[1:], ops)
	}
	if err != nil {
		return err
	}

	// run destroy (if applicable)
	if err := this.onDestroy(ctx, ops); err != nil {
		return err
	}

	return nil
}

// An internal representation of a command or subcommand and it's state before execution.
type command struct {
	Command

	fs       *flag.FlagSet
	args     []string
	showHelp bool
}

// onInit calls the [RunnableLifecycle] init routine if present on c.
func (c command) onInit(ctx context.Context, ops *ExecuteOptions) error {
	var err error

	if cmd, ok := c.Command.(RunnableLifecycle); ok {
		err = cmd.Initialize(ctx, c.args)
	}

	if errors.Is(err, ErrShowUsage) {
		_ = usage(c, ops)
		os.Exit(2)
	}

	return err
}

// run calls the [Runnable] run routine of c.
func (c command) run(ctx context.Context, ops *ExecuteOptions) error {
	err := c.Run(ctx, c.args)
	if errors.Is(err, ErrShowUsage) {
		_ = usage(c, ops)
		os.Exit(2)
	}

	return err
}

// onDestroy calls the [RunnableLifecycle] destroy routine if present on c.
func (c command) onDestroy(ctx context.Context, ops *ExecuteOptions) error {
	var err error

	if cmd, ok := c.Command.(RunnableLifecycle); ok {
		err = cmd.Destroy(ctx, c.args)
	}

	if errors.Is(err, ErrShowUsage) {
		_ = usage(c, ops)
		os.Exit(2)
	}

	return err
}

// buildCallStack builds a slice representing the command call stack. The first element in the slice is the root
// command and the last is the leaf command.
func buildCallStack(cmd Command, ops *ExecuteOptions) ([]command, error) {
	var stack []command

	var (
		args = ops.args
		err  error
	)

	for cmd != nil {
		this := command{
			Command: cmd,
			fs:      flag.NewFlagSet(cmd.Name(), flag.ContinueOnError),
		}

		// add help flags
		this.fs.BoolVar(&this.showHelp, "h", false, "show command help and usage information")
		this.fs.BoolVar(&this.showHelp, "help", false, "show command help and usage information")

		if c, ok := cmd.(FlagInitializer); ok {
			c.InitializeFlags(this.fs)
		}

		// bind environment variables
		if ops.bindEnv {
			if err := bindEnvironmentFlags(stack, this, ops); err != nil {
				return nil, err
			}
		}

		this.args, err = parseArgs(this, args, ops)
		if err != nil {
			return nil, err
		}

		args = this.args

		if len(args) == 0 {
			// if no subcommand name given, stop here
			cmd = nil
		} else if sub, ok := collectSubcommands(cmd)[args[0]]; ok {
			// if subcommand name given, continue
			args = args[1:]
			cmd = sub
		} else {
			// if arg given but it's not a subcommand name, stop here
			cmd = nil
		}

		stack = append(stack, this)
	}

	return stack, nil
}

// parseArgs processes args for the given command, returning the unparsed (remaining) arguments.
func parseArgs(cmd command, args []string, ops *ExecuteOptions) ([]string, error) {
	var fp flagParser = &getopt.PosixFlagSet{FlagSet: cmd.fs}

	if ops.nativeFlags {
		fp = cmd.fs
	}

	// interspersed args only possible for leaf commands
	interspersed := ops.interspersed
	if len(collectSubcommands(cmd.Command)) > 0 {
		interspersed = false
	}

	var processed []string

	for len(args) > 0 {
		if err := fp.Parse(args); err != nil {
			return nil, err
		}

		args = fp.Args()

		if !interspersed {
			return args, nil
		}

		if len(args) > 0 {
			processed = append(processed, args[0])
			args = args[1:]
		}
	}

	return processed, nil
}

// bindEnvironmentFlags sets flag values from matching environment variables.
func bindEnvironmentFlags(stack []command, cmd command, ops *ExecuteOptions) error {
	var components []string

	for _, c := range stack {
		components = append(components, c.Name())
	}

	components = append(components, cmd.Name())

	var flags []*flag.Flag
	cmd.fs.VisitAll(func(f *flag.Flag) {
		flags = append(flags, f)
	})

	for _, flag := range flags {
		variable := ops.bindEnvPrefix + formatEnvvar(append(components, flag.Name))

		if value, ok := os.LookupEnv(variable); ok {
			if err := flag.Value.Set(value); err != nil {
				return errors.Join(
					ErrEnvironmentBindFailure,
					fmt.Errorf("cmder: failed to set flag %s from variable %s", flag.Name, variable),
					err,
				)
			}
		}
	}

	return nil
}

// formatEnvvar generates an environment variable name which maps to the given flag path.
func formatEnvvar(flagpath []string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9]+")

	for i, v := range flagpath {
		flagpath[i] = strings.ToUpper(reg.ReplaceAllString(v, ""))
	}

	return strings.Join(flagpath, "_")
}

// helpRequested traverses the command stack and returns whether help text was requested with '-h' or '--help' flags,
// returning the leaf command from stack and true.
func helpRequested(stack []command) (*command, bool) {
	if len(stack) == 0 {
		return nil, false
	}

	for _, cmd := range stack {
		if cmd.showHelp {
			return &stack[len(stack)-1], true
		}
	}

	return nil, false
}
