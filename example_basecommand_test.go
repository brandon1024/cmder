package cmder_test

import (
	"context"
	"fmt"

	"github.com/brandon1024/cmder"
)

const BaseCommandExampleHelpText = `
'base-command' demonstrates how to build commands and subcommands with BaseCommand.
`

const BaseCommandExampleExamples = `
# broadcast hello to the world
base-command from cmder
`

type BaseCommandExample struct {
	cmder.BaseCommand
}

func (c *BaseCommandExample) Run(ctx context.Context, args []string) error {
	fmt.Printf("%s: %v\n", c.Name(), args)
	return nil
}

func ExampleBaseCommand() {
	cmd := &BaseCommandExample{
		BaseCommand: cmder.BaseCommand{
			CommandName: "base-command",
			Usage:       "base-command [<args>...]",
			ShortHelp:   "Simple demonstration of BaseCommand",
			Help:        BaseCommandExampleHelpText,
			Examples:    BaseCommandExampleExamples,
			Children: []cmder.Command{
				&cmder.BaseCommand{
					CommandName: "child",
					Usage:       "child [<args>...]",
					ShortHelp:   "A child command with simple behaviour",
					Help:        "I'm a simple subcommand with simple behaviour!",
					RunFunc: func(ctx context.Context, args []string) error {
						fmt.Printf("child: %v\n", args)
						return nil
					},
				},
			},
		},
	}

	args := []string{"child", "cmder"}
	if err := cmder.Execute(context.Background(), cmd, cmder.WithArgs(args)); err != nil {
		fmt.Printf("unexpected error occurred: %v", err)
	}

	// Output:
	// child: [cmder]
}
