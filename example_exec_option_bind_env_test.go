package cmder_test

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/brandon1024/cmder"
)

func ExampleWithEnvironmentBinding() {
	_ = os.Setenv("BINDENV_SHOW_FORMAT", "overidden-by-flag")
	_ = os.Setenv("BINDENV_SHOW_PAGECOUNT", "20")

	args := []string{"show", "--format=pretty"}

	ops := []cmder.ExecuteOption{
		cmder.WithArgs(args),
		cmder.WithEnvironmentBinding(),
	}

	if err := cmder.Execute(context.Background(), GetCommand(), ops...); err != nil {
		fmt.Printf("unexpected error occurred: %v", err)
	}
	// Output:
	// format: pretty
	// page-count: 20
}

const BindEnvHelpText = `
'bind-env' desmonstrates how cmder can be configured to bind environment variables to flags. Explicit command-line
arguments always take precedence over environment variables.
`

const BindEnvExamples = `
# print all default flag values
bind-env show
> default 10

# print flag values from environment
BINDENV_SHOW_FORMAT=pretty BINDENV_SHOW_PAGECOUNT=20 bind-env show --page-count=15
> pretty 15
`

func GetCommand() *cmder.BaseCommand {
	return &cmder.BaseCommand{
		CommandName: "bind-env",
		CommandDocumentation: cmder.CommandDocumentation{
			Usage:     "bind-env [subcommand] [flags]",
			ShortHelp: "Simple demonstration of binding environment variables to command flags.",
			Help:      BindEnvHelpText,
			Examples:  BindEnvExamples,
		},
		Children: []cmder.Command{GetShowCommand()},
	}
}

func GetShowCommand() *cmder.BaseCommand {
	return &cmder.BaseCommand{
		CommandName: "show",
		CommandDocumentation: cmder.CommandDocumentation{
			Usage:     `show [flags]`,
			ShortHelp: `Show flag values`,
			Help:      `'show' dumps flag values to stdout.`,
			Examples:  BindEnvExamples,
		},
		InitFlagsFunc: showFlags,
		RunFunc:       show,
	}
}

var (
	format string = "default"
	count  uint   = 10
)

func showFlags(fs *flag.FlagSet) {
	fs.StringVar(&format, "format", format, "output format (default, pretty)")
	fs.UintVar(&count, "page-count", count, "number of pages")
}

func show(ctx context.Context, args []string) error {
	switch format {
	case "default":
		fmt.Printf("%v %v\n", format, count)
	case "pretty":
		fmt.Printf("format: %v\npage-count: %v\n", format, count)
	default:
		return fmt.Errorf("illegal format: %s", format)
	}

	return nil
}
