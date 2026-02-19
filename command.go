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
//   - If you want to configure setup and teardown routines for a command, see [Initializer] and [Destroyer].
//   - If your command has subcommands, see [RootCommand].
//   - If your command has command-life flags and switches, see [FlagInitializer].
type Command interface {
	// All commands are [Runnable] and implement a Run routine.
	Runnable

	// Great tools come with great documentation. All commands need to provide documentation, which is used when
	// rendering usage and help texts.
	Documented

	// Name returns the name of this command or subcommand.
	Name() string
}

// Runnable is a fundamental interface implemented by commands. The Run routine is what carries out the work of your
// command.
//
// Concrete types can also implement [Initializer] or [Destroyer] to carry out any initialization and teardown
// necessary for your Run() routine.
type Runnable interface {
	// Run is the main body of your command executed by [Execute].
	//
	// The given [context.Context] is derived from the context provided to Execute() and is cancelled when Execute()
	// returns. Use this context to cleanup resources.
	//
	// The second argument is the list of command-line arguments and switches that remain after parsing flags.
	Run(context.Context, []string) error
}

// Initializer may be implemented by commands that need to do some work before the [Runnable] Run() routine is invoked.
//
// See [Execute] for more details on the lifecycle of command execution.
type Initializer interface {
	// Initialize carries out any initialization needed for this [Command]. Errors returned by Initialize will abort
	// execution of the command lifecycle (Run()/Destroy() of this command and parent command(s)).
	Initialize(context.Context, []string) error
}

// Destroyer may be implemented by commands that need to do some work after the [Runnable] Run() routine is invoked.
//
// See [Execute] for more details on the lifecycle of command execution.
type Destroyer interface {
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
	// UsageLine returns the usage line for your command. This is akin to the "SYNOPSIS" section you would typically
	// find in a man page. Generally, usage lines have a well accepted format:
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

	// ShortHelpText returns a short-and-sweet one-line description of your command. This is akin to the "NAME" section
	// you would typically find in a man page. This is mainly used to summarize available subcommands.
	//
	// Here are a few examples:
	//
	//	get image registry and namespace information
	//	download installation files for an application version
	//	display the contents of a file in a terminal
	ShortHelpText() string

	// HelpText returns longer usage and help information for your users about this subcommand. Here you can describe
	// the behaviour of your command, summarize usage of certain flags and arguments and provide hints on where to find
	// additional information. This is akin to the "DESCRIPTION" section you would typically find in a man page.
	//
	// For a better viewing experience in terminals, consider maintaining consistent line length limits (120 is a good
	// target).
	HelpText() string

	// ExampleText returns motivating usage examples for your command.
	ExampleText() string
}

// HiddenCommand is implemented by commands which are not user facing. Hidden commands are not displayed in help texts.
type HiddenCommand interface {
	// Hidden returns a flag indicating whether to mark this command as hidden, preventing it from being rendered in
	// help output.
	Hidden() bool
}

// Compile-time checks.
var (
	_ Command         = &BaseCommand{}
	_ Initializer     = &BaseCommand{}
	_ Destroyer       = &BaseCommand{}
	_ RootCommand     = &BaseCommand{}
	_ FlagInitializer = &BaseCommand{}
	_ Documented      = &CommandDocumentation{}
	_ HiddenCommand   = &CommandDocumentation{}
)

// CommandDocumentation implements [Documented] and can be embdded in command types to reduce boilerplate.
type CommandDocumentation struct {
	// The usage line. See UsageLine() in [Documented].
	Usage string

	// The short help line. See ShortHelpText() in [Documented].
	ShortHelp string

	// Documentation for your command. See HelpText() in [Documented].
	Help string

	// Usage examples for your command. See ExampleText() in [Documented].
	Examples string

	// Whether this command is hidden in help and usage texts. See Hidden() in [HiddenCommand].
	IsHidden bool
}

// UsageLine returns [CommandDocumentation] Usage.
//
// See [Documented].
func (d CommandDocumentation) UsageLine() string {
	return d.Usage
}

// ShortHelpText returns [CommandDocumentation] ShortHelp.
//
// See [Documented].
func (d CommandDocumentation) ShortHelpText() string {
	return d.ShortHelp
}

// HelpText returns [CommandDocumentation] Help.
//
// See [Documented].
func (d CommandDocumentation) HelpText() string {
	return d.Help
}

// ExampleText returns [CommandDocumentation] Examples.
//
// See [Documented].
func (d CommandDocumentation) ExampleText() string {
	return d.Examples
}

// Hidden returns [CommandDocumentation] Hidden.
//
// See [HiddenCommand].
func (d CommandDocumentation) Hidden() bool {
	return d.IsHidden
}

// BaseCommand is an implementation of the [Command], [Initializer], [Destroyer], [RootCommand] and [FlagInitializer]
// interfaces and may be embedded in your command types to reduce boilerplate.
type BaseCommand struct {
	CommandDocumentation

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
// See [Initializer].
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
// See [Destroyer].
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
