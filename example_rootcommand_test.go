package cmder_test

import (
	"context"
	"fmt"

	"github.com/brandon1024/cmder"
)

func ExampleRootCommand() {
	cmd := &ParentCommand{
		subcommands: []cmder.Command{&ChildCommand{}},
	}

	args := []string{"child", "hello-world"}
	if err := cmder.Execute(context.Background(), cmd, cmder.WithArgs(args)); err != nil {
		fmt.Printf("unexpected error occurred: %v", err)
	}
	// Output:
	// parent: init [child hello-world]
	// child: init [hello-world]
	// child: run [hello-world]
	// child: destroy [hello-world]
	// parent: destroy [child hello-world]
}

// === PARENT COMMAND ===

const ParentCommandUsageLine = `parent [<subcommand>] [<args>...]`

const ParentCommandShortHelpText = `Example of parent command`

const ParentCommandHelpText = `
'parent' demonstrates an example of a command with subcommands. When executed without any arguments, the parent's Run
routine is executed, but if the child subcommand is provided the child subcommand Run routine will be exeuted instead.

The parent implements RootCommand indicating that it is a root command with runnable subcommands. The child does not
implement RootCommand, indicating it is a leaf command.

In this example, the parent and child commands implement RunnableLifecycle, and their respective init/destroy routines
are invoked in this order:

  1. parent Initialize
  2. child Initialize
  3. child Run
  4. child Destroy
  5. parent Destroy
`

const ParentCommandExamples = `
# run the parent 'Run' routine
parent

# run the child 'Run' routine
parent child

# run with some additional args
parent hello-world
parent child hello-world
`

type ParentCommand struct {
	subcommands []cmder.Command
}

func (c *ParentCommand) Name() string {
	return "parent"
}

func (c *ParentCommand) Initialize(ctx context.Context, args []string) error {
	fmt.Printf("%s: init %v\n", c.Name(), args)
	return nil
}

func (c *ParentCommand) Run(ctx context.Context, args []string) error {
	fmt.Printf("%s: run %v\n", c.Name(), args)
	return nil
}

func (c *ParentCommand) Destroy(ctx context.Context, args []string) error {
	fmt.Printf("%s: destroy %v\n", c.Name(), args)
	return nil
}

func (c *ParentCommand) UsageLine() string {
	return ParentCommandUsageLine
}

func (c *ParentCommand) ShortHelpText() string {
	return ParentCommandShortHelpText
}

func (c *ParentCommand) HelpText() string {
	return ParentCommandHelpText
}

func (c *ParentCommand) ExampleText() string {
	return ParentCommandExamples
}

func (c *ParentCommand) Subcommands() []cmder.Command {
	return c.subcommands
}

// === CHILD COMMAND ===

const ChildCommandUsageLine = `child [<args>...]`

const ChildCommandShortHelpText = `Example of child command`

const ChildCommandHelpText = `
'child' is the subcommand of 'parent'.
`

const ChildCommandExamples = `
# run the child 'Run' routine
parent child

# run with some additional args
parent child hello-world
`

type ChildCommand struct{}

func (c *ChildCommand) Name() string {
	return "child"
}

func (c *ChildCommand) Initialize(ctx context.Context, args []string) error {
	fmt.Printf("%s: init %v\n", c.Name(), args)
	return nil
}

func (c *ChildCommand) Run(ctx context.Context, args []string) error {
	fmt.Printf("%s: run %v\n", c.Name(), args)
	return nil
}

func (c *ChildCommand) Destroy(ctx context.Context, args []string) error {
	fmt.Printf("%s: destroy %v\n", c.Name(), args)
	return nil
}

func (c *ChildCommand) UsageLine() string {
	return ChildCommandUsageLine
}

func (c *ChildCommand) ShortHelpText() string {
	return ChildCommandShortHelpText
}

func (c *ChildCommand) HelpText() string {
	return ChildCommandHelpText
}

func (c *ChildCommand) ExampleText() string {
	return ChildCommandExamples
}
