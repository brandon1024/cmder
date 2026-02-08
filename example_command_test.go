package cmder_test

import (
	"context"
	"fmt"

	"github.com/brandon1024/cmder"
)

func ExampleCommand() {
	cmd := &HelloWorldCommand{}

	args := []string{"from", "cmder"}

	if err := cmder.Execute(context.Background(), cmd, cmder.WithArgs(args)); err != nil {
		fmt.Printf("unexpected error occurred: %v", err)
	}
	// Output:
	// hello-world: [from cmder]
}

const HelloWorldCommandUsageLine = `hello-world [<args>...]`

const HelloWorldCommandShortHelpText = `Simple demonstration of cmder`

const HelloWorldCommandHelpText = `
'hello-world' demonstrates the simplest usage of cmder. This example defines a single command 'hello-world' that
implements the Runnable and Documented interfaces.
`

const HelloWorldCommandExamples = `
# broadcast hello to the world
hello-world from cmder
`

type HelloWorldCommand struct{}

func (c *HelloWorldCommand) Name() string {
	return "hello-world"
}

func (c *HelloWorldCommand) Run(ctx context.Context, args []string) error {
	fmt.Printf("%s: %v\n", c.Name(), args)
	return nil
}

func (c *HelloWorldCommand) UsageLine() string {
	return HelloWorldCommandUsageLine
}

func (c *HelloWorldCommand) ShortHelpText() string {
	return HelloWorldCommandShortHelpText
}

func (c *HelloWorldCommand) HelpText() string {
	return HelloWorldCommandHelpText
}

func (c *HelloWorldCommand) ExampleText() string {
	return HelloWorldCommandExamples
}
