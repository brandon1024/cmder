package main

import (
	"context"
	"fmt"
)

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

func (c *ChildCommand) Hidden() bool {
	return false
}
