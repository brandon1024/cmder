package cmder_test

import (
	"context"
	"flag"
	"fmt"

	"github.com/brandon1024/cmder"
)

// === PARENT COMMAND ===

const ParentFlagsCommandUsageLine = `parent-flags [<subcommand>] [<args>...]`

const ParentFlagsCommandShortHelpText = `Example of parent command with flags`

const ParentFlagsCommandHelpText = `
'parent-flags' demonstrates an example of a command with subcommands and flags. Commands can define flags by
implementing the FlagInitializer interface. During execution, flags are parsed at each level in the command tree and
can be accessed by the command's Run routine.
`

const ParentFlagsCommandExamples = `
# run the parent 'Run' routine with flags
parent-flags --option string

# run the child 'Run' routine with flags
parent-flags child-flags --count 1

# run with some additional args
parent-flags --option string hello-world
parent-flags --option string child-flags --count 1 hello-world
`

type ParentFlagsCommand struct {
	cmder.BaseCommand

	option string
}

func (c *ParentFlagsCommand) Initialize(ctx context.Context, args []string) error {
	fmt.Printf("%s: init [%s] %v\n", c.Name(), c.option, args)
	return nil
}

func (c *ParentFlagsCommand) Run(ctx context.Context, args []string) error {
	fmt.Printf("%s: run [%s] %v\n", c.Name(), c.option, args)
	return nil
}

func (c *ParentFlagsCommand) Destroy(ctx context.Context, args []string) error {
	fmt.Printf("%s: destroy [%s] %v\n", c.Name(), c.option, args)
	return nil
}

func (c *ParentFlagsCommand) InitializeFlags(fs *flag.FlagSet) {
	fs.StringVar(&c.option, "option", "", "parent command option string argument")
}

func NewParentFlagsCommand() *ParentFlagsCommand {
	return &ParentFlagsCommand{
		BaseCommand: cmder.BaseCommand{
			CommandName: "parent-flags",
			Usage:       ParentFlagsCommandUsageLine,
			ShortHelp:   ParentFlagsCommandShortHelpText,
			Help:        ParentFlagsCommandHelpText,
			Examples:    ParentFlagsCommandExamples,
			Children: []cmder.Command{
				NewChildFlagsCommand(),
			},
		},
	}
}

// === CHILD COMMAND ===

const ChildFlagsCommandUsageLine = `child-flags [<args>...]`

const ChildFlagsCommandShortHelpText = `Example of child command`

const ChildFlagsCommandHelpText = `
'child-flags' is the subcommand of 'parent-flags'.
`

const ChildFlagsCommandExamples = `
# run the child 'Run' routine with flags
parent-flags child-flags --count 1

# run with some additional args
parent-flags --option string child-flags --count 1 hello-world
`

type ChildFlagsCommand struct {
	cmder.BaseCommand

	count int
}

func (c *ChildFlagsCommand) Initialize(ctx context.Context, args []string) error {
	fmt.Printf("%s: init [%d] %v\n", c.Name(), c.count, args)
	return nil
}

func (c *ChildFlagsCommand) Run(ctx context.Context, args []string) error {
	fmt.Printf("%s: run [%d] %v\n", c.Name(), c.count, args)
	return nil
}

func (c *ChildFlagsCommand) Destroy(ctx context.Context, args []string) error {
	fmt.Printf("%s: destroy [%d] %v\n", c.Name(), c.count, args)
	return nil
}

func (c *ChildFlagsCommand) InitializeFlags(fs *flag.FlagSet) {
	fs.IntVar(&c.count, "count", 0, "child command count integer argument")
}

func NewChildFlagsCommand() *ChildFlagsCommand {
	return &ChildFlagsCommand{
		BaseCommand: cmder.BaseCommand{
			CommandName: "child-flags",
			Usage:       ChildFlagsCommandUsageLine,
			ShortHelp:   ChildFlagsCommandShortHelpText,
			Help:        ChildFlagsCommandHelpText,
			Examples:    ChildFlagsCommandExamples,
		},
	}
}

// === EXAMPLE ===

func ExampleFlagInitializer() {
	cmd := NewParentFlagsCommand()

	args := []string{"--option", "test=1", "child-flags", "--count", "1", "hello-world"}
	if err := cmder.Execute(context.Background(), cmd, cmder.WithArgs(args)); err != nil {
		fmt.Printf("unexpected error occurred: %v", err)
	}

	// Output:
	// parent-flags: init [test=1] [child-flags --count 1 hello-world]
	// child-flags: init [1] [hello-world]
	// child-flags: run [1] [hello-world]
	// child-flags: destroy [1] [hello-world]
	// parent-flags: destroy [test=1] [child-flags --count 1 hello-world]
}
