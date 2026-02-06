package cmder

import (
	"context"
	"flag"
)

// Command is the fundamental interface implemented by types that are runnable commands or subcommands. Commands can
// be executed with [Execute].
//
// Concrete types can implement additional interfaces to configure additional behaviour, like setup/teardown routines,
// subcommands, command-line flags, and other behaviour:
//
//   - If you want to configure setup and teardown routines for a command, see [RunnableLifecycle].
//   - If your command has subcommands, see [RootCommand].
//   - If your command has command-life flags and switches, see [FlagInitializer].
type Command interface {
	// All commands are [Runnable] and implement a Run routine.
	Runnable

	// Great tools come with great documentation. All commands need to provide documentation, which is used when
	// rendering usage and help texts.
	Documented

	// Name returns the name of this command.
	Name() string
}

// Runnable is a fundamental interface implemented by commands. The Run routine is what carries out the work of your
// command.
//
// Concrete types can also implement [RunnableLifecycle] to carry out any initialization and teardown necessary for your
// Run() routine.
type Runnable interface {
	// Run is the main body of your command executed by [Execute].
	//
	// The given [context.Context] is derived from the context provided to Execute() and is cancelled when Execute()
	// returns. Use this context to cleanup resources.
	//
	// The second argument is the list of command-line arguments and switches that remain after parsing flags.
	Run(context.Context, []string) error
}

// RunnableLifecycle may be implemented by commands that need to do some work before and after the [Runnable] Run()
// routine is invoked.
//
// When executing subcommands, the Initialize() and Destroy() routines of parent commands are also invoked. For
// instance, if executing subcommand 'child' of command 'parent', lifecycle routines are invoked in this order:
//
//  1. parent: Initialize()
//  2. child: Initialize()
//  3. child: Run()
//  4. child: Destroy()
//  5. parent: Destroy()
//
// When executing subcommands, the arguments provided to the Initialize() and Destroy() routines of parent commands will
// include the unprocessed args and flags of child commands. For example:
//
//	$ parent --option test child --count 1 arg-1 arg-2
//
// will execute lifecycle routines with the arguments:
//
//  1. parent: Initialize [child --count 1 arg-1 arg-2]
//  2. child: Initialize  [arg-1 arg-2]
//  3. child: Run         [arg-1 arg-2]
//  4. child: Destroy     [arg-1 arg-2]
//  5. parent: Destroy    [child --count 1 arg-1 arg-2]
type RunnableLifecycle interface {
	// Initialize carries out any initialization needed for this [Command]. Errors returned by Initialize will abort
	// execution of the command lifecycle (Run()/Destroy() of this command and parent command(s)).
	Initialize(context.Context, []string) error

	// Destroy carries out any teardown needed for this [Command]. Errors returned by Destroy will abort execution of
	// the command lifecycle (Destroy of this command and parent command(s)).
	Destroy(context.Context, []string) error
}

// RootCommand may be implemented by commands that have subcommands.
type RootCommand interface {
	// Subcommands returns a slice of subcommands of this RootCommand. May return nil or an empty slice to treat this
	// command as a leaf command.
	Subcommands() []Command
}

// Documented is implemented by all commands and provides help and usage information for your users.
type Documented interface {
	// UsageLine returns the usage line for your command. Generally, usage lines have a well accepted format:
	//
	//	- [ ] identifies an optional argument or flag. Arguments that are not enclosed in brackets are required.
	//	- ... identifies arguments or flags that can be provided more than once.
	//	-  |  identifies mutually exclusive arguments or flags.
	//	- ( ) identifies groups of flags or arguments that are required together.
	//	- < > identifies argument(s) or flag(s).
	//
	// Here are a few examples:
	//
	//	git add [<options>] [--] <pathspec>...
	//	kubectl get [(-o|--output=)json|yaml|wide] (TYPE[.VERSION][.GROUP] [NAME | -l label] | TYPE[.VERSION][.GROUP]/NAME ...) [flags] [options]
	//	crane index filter [flags]
	UsageLine() string

	// ShortHelpText returns a short-and-sweet one-line description of your command. This is mainly used to summarize
	// available subcommands.
	ShortHelpText() string

	// HelpText returns long usage and help information for your users about this subcommand. Here you can describe the
	// behaviour of your command, summarize usage of certain flags and arguments and provide hints on where to find
	// additional information. This is akin to the "DESCRIPTION" section you would typically find in a man page.
	//
	// For a better viewing experience in terminals, consider maintaining consistent line length limits (120 is a good
	// target).
	HelpText() string

	// ExampleText returns motivating usage examples for your command.
	ExampleText() string

	// Hidden returns a flag indicating whether to mark this command as hidden, preventing it from being rendered in
	// help output.
	Hidden() bool
}

// Compile-time checks.
var (
	_ Command           = &BaseCommand{}
	_ RunnableLifecycle = &BaseCommand{}
	_ RootCommand       = &BaseCommand{}
	_ FlagInitializer   = &BaseCommand{}
)

// BaseCommand is an implementation of the [Command], [RunnableLifecycle], [RootCommand] and [FlagInitializer]
// interfaces and may be embedded in your command types to reduce boilerplate.
type BaseCommand struct {
	// The command name. See Name() in [Command].
	CommandName string

	// Optional function invoked by the default InitializeFlags() function.
	InitFlagsFunc func(*flag.FlagSet)

	// Optional function invoked by the default Initialize() function.
	InitFunc func(context.Context, []string) error

	// Optional function invoked by the default Run() function.
	RunFunc func(context.Context, []string) error

	// Optional function invoked by the default Destroy() function.
	DestroyFunc func(context.Context, []string) error

	// Subcommands for this command, if applicable. See [RootCommand].
	Children []Command

	// The usage line. See UsageLine() in [Documented].
	Usage string

	// The short help line. See ShortHelpText() in [Documented].
	ShortHelp string

	// Documentation for your command. See HelpText() in [Documented].
	Help string

	// Usage examples for your command. See ExampleText() in [Documented].
	Examples string

	// Whether this command is hidden in help and usage texts. See Hidden() in [Documented].
	IsHidden bool
}

// Name returns [BaseCommand] CommandName.
//
// See [Command].
func (c BaseCommand) Name() string {
	return c.CommandName
}

// InitializeFlags runs [BaseCommand] InitFlagsFunc, if not nil.
//
// See [FlagInitializer].
func (c BaseCommand) InitializeFlags(fs *flag.FlagSet) {
	if c.InitFlagsFunc != nil {
		c.InitFlagsFunc(fs)
	}
}

// Initialize runs [BaseCommand] InitFunc, if not nil.
//
// See [RunnableLifecycle].
func (c BaseCommand) Initialize(ctx context.Context, args []string) error {
	if c.InitFunc != nil {
		return c.InitFunc(ctx, args)
	}

	return nil
}

// Run runs [BaseCommand] RunFunc, if not nil.
//
// See [Runnable].
func (c BaseCommand) Run(ctx context.Context, args []string) error {
	if c.RunFunc != nil {
		return c.RunFunc(ctx, args)
	}

	return nil
}

// Destroy runs [BaseCommand] DestroyFunc, if not nil.
//
// See [RunnableLifecycle].
func (c BaseCommand) Destroy(ctx context.Context, args []string) error {
	if c.DestroyFunc != nil {
		return c.DestroyFunc(ctx, args)
	}

	return nil
}

// Subcommands returns [BaseCommand] Children.
//
// See [RootCommand].
func (c BaseCommand) Subcommands() []Command {
	return c.Children
}

// UsageLine returns [BaseCommand] Usage.
//
// See [Documented].
func (c BaseCommand) UsageLine() string {
	return c.Usage
}

// ShortHelpText returns [BaseCommand] ShortHelp.
//
// See [Documented].
func (c BaseCommand) ShortHelpText() string {
	return c.ShortHelp
}

// HelpText returns [BaseCommand] Help.
//
// See [Documented].
func (c BaseCommand) HelpText() string {
	return c.Help
}

// ExampleText returns [BaseCommand] Examples.
//
// See [Documented].
func (c BaseCommand) ExampleText() string {
	return c.Examples
}

// Hidden returns [BaseCommand] Hidden.
//
// See [Documented].
func (c BaseCommand) Hidden() bool {
	return c.IsHidden
}
