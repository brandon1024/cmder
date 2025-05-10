package cmder

import (
	"context"
	"errors"
	"flag"
	"os"
)

var (
	// Returned when a [Command] provided to [Execute] is illegal.
	ErrIllegalCommandConfiguration = errors.New("illegal command configuration")

	// Returned when an [ExecuteOption] provided to [Execute] is illegal.
	ErrIllegalExecuteOptions = errors.New("illegal command execution option")
)

// A no-op implementation of [RunnableLifecycle] used internally by [Execute].
type noOpRunnableLifecycle struct{}

func (i *noOpRunnableLifecycle) Initialize(context.Context, []string) error {
	return nil
}

func (i *noOpRunnableLifecycle) Destroy(context.Context, []string) error {
	return nil
}

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
// immediately and the error is returned. For example, returning an error from Run() will prevent execution of Destroy()
// of the current command and any parents.
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
// command-line flags. Each command/subcommand is given a unique [flag.FlagSet].
//
// If a command does not define help flags (-h, --help), Execute will register help flags and will display command usage
// text to [UsageOutputWriter] if those flags are provided at the command line.
func Execute(ctx context.Context, cmd Command, op ...ExecuteOption) error {
	// setup context
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// prepare executor options
	ops := &ExecuteOptions{
		args:  os.Args[1:],
		flags: flag.NewFlagSet(cmd.Name(), flag.ContinueOnError),
	}
	for _, f := range op {
		f(ops)
	}

	// do some checks
	if ops.flags == nil {
		return errors.Join(
			ErrIllegalExecuteOptions,
			errors.New("command flag set cannot be nil"),
		)
	}

	// initialize flags
	if c, ok := cmd.(FlagInitializer); ok {
		c.InitializeFlags(ops.flags)
	}

	// if no help flags are defined, register them here
	var showHelp bool
	if ops.flags.Lookup("help") == nil && ops.flags.Lookup("h") == nil {
		ops.flags.BoolVar(&showHelp, "help", false, "show command help and usage text")
		ops.flags.BoolVar(&showHelp, "h", false, "show command help and usage text")
	}

	// parse flags
	if err := ops.flags.Parse(ops.args); err != nil {
		return err
	}

	// if help flag found, display help pages
	if showHelp {
		return RenderUsage(cmd, ops.flags)
	}

	args := ops.flags.Args()

	// if cmd has a lifecycle, use it, otherwise fallback to a no-op lifecycle
	var lifecycle RunnableLifecycle = &noOpRunnableLifecycle{}
	if l, ok := cmd.(RunnableLifecycle); ok {
		lifecycle = l
	}

	// run command initialization
	if err := lifecycle.Initialize(ctx, args); err != nil {
		return err
	}

	if c, ok := cmd.(RootCommand); ok {
		// if this is a root command, run appropriate subcommand if given
		subcommands := collectSubcommands(c)

		// if no args given, run this command
		if len(args) == 0 {
			if err := cmd.Run(ctx, args); err != nil {
				return err
			}
		} else {
			if sub, ok := subcommands[args[0]]; !ok {
				// if first arg isn't a subcommand, run this command
				if err := cmd.Run(ctx, args); err != nil {
					return err
				}
			} else {
				// else, run subcommand
				if err := Execute(ctx, sub, WithArgs(args[1:])); err != nil {
					return err
				}
			}
		}
	} else {
		// if this is a leaf command, run it
		if err := cmd.Run(ctx, args); err != nil {
			return err
		}
	}

	// run command teardown
	if err := lifecycle.Destroy(ctx, args); err != nil {
		return err
	}

	return nil
}

// collectSubcommands collects the immediate subcommands of the given [RootCommand] into a map keyed by the command
// [Command] Name().
func collectSubcommands(cmd RootCommand) map[string]Command {
	subcommands := map[string]Command{}
	for _, subcommand := range cmd.Subcommands() {
		subcommands[subcommand.Name()] = subcommand
	}

	return subcommands
}
