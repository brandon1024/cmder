package cmder_test

import (
	"context"
	"fmt"

	"github.com/brandon1024/cmder"
)

const LifecycleCommandUsageLine = `lifecycle [<args>...]`

const LifecycleCommandShortHelpText = `Example command with lifecycle routines`

const LifecycleCommandHelpText = `
'lifecycle' demonstrates a command that implements the RunnableLifecycle interface, defining initialization and
destroy routines.
`

const LifecycleCommandExamples = `
# demonstrate initialization and teardown
lifecycle
`

type LifecycleCommand struct{}

func (c *LifecycleCommand) Name() string {
	return "lifecycle"
}

func (c *LifecycleCommand) Initialize(ctx context.Context, args []string) error {
	fmt.Println("lifecycle: initializing")
	return nil
}

func (c *LifecycleCommand) Run(ctx context.Context, args []string) error {
	fmt.Println("lifecycle: running")
	return nil
}

func (c *LifecycleCommand) Destroy(ctx context.Context, args []string) error {
	fmt.Println("lifecycle: shutting down")
	return nil
}

func (c *LifecycleCommand) UsageLine() string {
	return LifecycleCommandUsageLine
}

func (c *LifecycleCommand) ShortHelpText() string {
	return LifecycleCommandShortHelpText
}

func (c *LifecycleCommand) HelpText() string {
	return LifecycleCommandHelpText
}

func (c *LifecycleCommand) ExampleText() string {
	return LifecycleCommandExamples
}

func (c *LifecycleCommand) Hidden() bool {
	return false
}

func ExampleRunnableLifecycle() {
	args := []string{}

	cmd := &LifecycleCommand{}

	if err := cmder.Execute(context.Background(), cmd, cmder.WithArgs(args)); err != nil {
		fmt.Printf("unexpected error occurred: %v", err)
	}

	// Output:
	// lifecycle: initializing
	// lifecycle: running
	// lifecycle: shutting down
}
