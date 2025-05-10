package main

import (
	"context"
	"fmt"
)

const WorldCommandUsageLine = `hello world [<args>...]`

const WorldCommandShortHelpText = `Hello world example subcommand`

const WorldCommandHelpText = `
'world' is the subcommand of 'hello'.
`

const WorldCommandExamples = `
# run the child 'Run' routine
hello world

# run with some additional args
hello world from cmder
`

type WorldCommand struct{}

func (c *WorldCommand) Name() string {
	return "world"
}

func (c *WorldCommand) Initialize(ctx context.Context, args []string) error {
	fmt.Printf("%s: init %v\n", c.Name(), args)
	return nil
}

func (c *WorldCommand) Run(ctx context.Context, args []string) error {
	fmt.Printf("%s: run %v\n", c.Name(), args)
	return nil
}

func (c *WorldCommand) Destroy(ctx context.Context, args []string) error {
	fmt.Printf("%s: destroy %v\n", c.Name(), args)
	return nil
}

func (c *WorldCommand) UsageLine() string {
	return WorldCommandUsageLine
}

func (c *WorldCommand) ShortHelpText() string {
	return WorldCommandShortHelpText
}

func (c *WorldCommand) HelpText() string {
	return WorldCommandHelpText
}

func (c *WorldCommand) ExampleText() string {
	return WorldCommandExamples
}

func (c *WorldCommand) Hidden() bool {
	return false
}
