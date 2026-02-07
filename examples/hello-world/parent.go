package main

import (
	"context"
	"fmt"

	"github.com/brandon1024/cmder"
)

const HelloCommandUsageLine = `hello [<subcommand>] [<args>...]`

const HelloCommandShortHelpText = `Hello world example command`

const HelloCommandHelpText = `
'hello' demonstrates an example of a command with subcommands. When executed without any arguments, 'hello' executes
the Run routine for the root command, but if the 'world' subcommand is provided the 'world' subcommand Run routine will
be exeuted instead.

'hello' implements RootCommand indicating that it is a root command with runnable subcommands. 'world' does not
implement RootCommand, indicating it is a leaf command.

In this example, the 'hello' and 'world' commands implement RunnableLifecycle, and their respective init/destroy
routines are invoked in this order:

  1. hello Initialize
  2. world Initialize
  3. world Run
  4. world Destroy
  5. hello Destroy
`

const HelloCommandExamples = `
# run the parent 'Run' routine
hello

# run the child 'Run' routine
hello world

# run with some additional args
hello from cmder
hello world from cmder
`

type HelloCommand struct {
	subcommands []cmder.Command
}

func (c *HelloCommand) Name() string {
	return "hello"
}

func (c *HelloCommand) Initialize(ctx context.Context, args []string) error {
	fmt.Printf("%s: init %v\n", c.Name(), args)
	return nil
}

func (c *HelloCommand) Run(ctx context.Context, args []string) error {
	fmt.Printf("%s: run %v\n", c.Name(), args)
	return nil
}

func (c *HelloCommand) Destroy(ctx context.Context, args []string) error {
	fmt.Printf("%s: destroy %v\n", c.Name(), args)
	return nil
}

func (c *HelloCommand) UsageLine() string {
	return HelloCommandUsageLine
}

func (c *HelloCommand) ShortHelpText() string {
	return HelloCommandShortHelpText
}

func (c *HelloCommand) HelpText() string {
	return HelloCommandHelpText
}

func (c *HelloCommand) ExampleText() string {
	return HelloCommandExamples
}

func (c *HelloCommand) Subcommands() []cmder.Command {
	return c.subcommands
}
