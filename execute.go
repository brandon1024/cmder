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
func Execute(ctx context.Context, cmd Command, op ...ExecuteOption) error {
	// do some checks
	if cmd == nil {
		return errors.Join(ErrIllegalCommandConfiguration, errors.New("command cannot be nil"))
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

	this := stack[0]

	// run init, if applicable
	if l, ok := this.cmd.(RunnableLifecycle); ok {
		if err := l.Initialize(ctx, this.fs.Args()); err != nil {
			return err
		}
	}

	// if this is a leaf, run Run(), otherwise recurse
	var err error
	if len(stack) == 1 {
		err = this.cmd.Run(ctx, this.fs.Args())
	} else {
		err = execute(ctx, stack[1:])
	}

	if err != nil {
		return err
	}

	// run destroy, if applicable
	if l, ok := this.cmd.(RunnableLifecycle); ok {
		if err := l.Destroy(ctx, this.fs.Args()); err != nil {
			return err
		}
	}

	return nil
}

// collectSubcommands collects the immediate subcommands of the given [Command] into a map keyed by the command
// [Command] Name(). Returns an empty map if the command is not a [RootCommand].
func collectSubcommands(cmd Command) map[string]Command {
	subcommands := map[string]Command{}

	c, ok := cmd.(RootCommand)
	if !ok {
		return subcommands
	}

	for _, subcommand := range c.Subcommands() {
		subcommands[subcommand.Name()] = subcommand
	}

	return subcommands
}

// An internal representation of a command or subcommand and it's state before execution.
type command struct {
	cmd Command
	fs  *flag.FlagSet
}

// buildCallStack builds a slice representing the command call stack. The first element in the slice is the root
// command and the last is the leaf command.
func buildCallStack(cmd Command, args []string) ([]command, error) {
	var stack []command

	if cmd == nil {
		return nil, errors.Join(ErrIllegalCommandConfiguration, errors.New("command cannot be nil"))
	}

	for cmd != nil {
		this := command{
			cmd: cmd,
			fs:  flag.NewFlagSet(cmd.Name(), flag.ContinueOnError),
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
