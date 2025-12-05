package cmder

import (
	"context"
	"errors"
	"flag"
	"os"
)

var (
	// Returned when a [Command] provided to [Execute] is illegal.
	ErrIllegalCommandConfiguration = errors.New("cmder: illegal command configuration")

	// Returned when an [ExecuteOption] provided to [Execute] is illegal.
	ErrIllegalExecuteOptions = errors.New("cmder: illegal command execution option")
)

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
// Execute may return [ErrIllegalCommandConfiguration] or [ErrIllegalExecuteOptions] if a command is misconfigured or
// options are invalid.
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
// # Usage and Help Texts
//
// Whenever the user provides the '-h' or '--help' flag at the command line, [Execute] will display command usage and
// exit. The format of the help text can be adjusted by configuring [UsageTemplate]. By default, usage information will
// be written to stderr, but this can be adjusted by setting [UsageOutputWriter].
func Execute(ctx context.Context, cmd Command, op ...ExecuteOption) error {
	// do some checks
	if cmd == nil {
		return errors.Join(ErrIllegalCommandConfiguration, errors.New("cmder: command cannot be nil"))
	}

	// prepare executor options
	ops := &ExecuteOptions{
		args: os.Args[1:],
	}
	for _, f := range op {
		f(ops)
	}

	// build a stack of command invocations
	stack, err := buildCallStack(cmd, ops.args)
	if err != nil {
		return err
	}

	// if help was requested, display and exit
	if cmd, ok := helpRequested(stack); ok {
		return usage(*cmd)
	}

	return execute(ctx, stack)
}

// execute traverses the command stack recursively executing the lifecycle routines at each level.
func execute(ctx context.Context, stack []command) error {
	if len(stack) == 0 {
		return nil
	}

	// setup context
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var (
		this = stack[0]
		args = this.fs.Args()
		err  error
	)

	// run init (if applicable)
	if err := this.initializeFn(ctx, args); err != nil {
		return err
	}

	// if this is a leaf, run Run(), otherwise recurse
	if len(stack) == 1 {
		err = this.Run(ctx, args)
	} else {
		err = execute(ctx, stack[1:])
	}
	if err != nil {
		return err
	}

	// run destroy (if applicable)
	if err := this.destroyFn(ctx, args); err != nil {
		return err
	}

	return nil
}

// An internal representation of a command or subcommand and it's state before execution.
type command struct {
	Command

	fs           *flag.FlagSet
	initializeFn func(context.Context, []string) error
	destroyFn    func(context.Context, []string) error
	showHelp     bool
}

// buildCallStack builds a slice representing the command call stack. The first element in the slice is the root
// command and the last is the leaf command.
func buildCallStack(cmd Command, args []string) ([]command, error) {
	var stack []command

	for cmd != nil {
		this := command{
			Command:      cmd,
			fs:           flag.NewFlagSet(cmd.Name(), flag.ContinueOnError),
			initializeFn: func(context.Context, []string) error { return nil },
			destroyFn:    func(context.Context, []string) error { return nil },
		}

		// add help flags
		this.fs.BoolVar(&this.showHelp, "h", false, "show command help and usage information")
		this.fs.BoolVar(&this.showHelp, "help", false, "show command help and usage information")

		if l, ok := cmd.(RunnableLifecycle); ok {
			this.initializeFn = l.Initialize
			this.destroyFn = l.Destroy
		}

		if c, ok := cmd.(FlagInitializer); ok {
			c.InitializeFlags(this.fs)
		}

		if err := this.fs.Parse(args); err != nil {
			return nil, err
		}

		args = this.fs.Args()

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
