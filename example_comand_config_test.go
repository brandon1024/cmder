package cmder_test

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/brandon1024/cmder"
	"github.com/brandon1024/cmder/getopt"
)

// This example demonstrates how to configure your application from a configuration file, environment variables, and
// flags.
func ExampleCommand_config() {
	// lowest precedence
	json.Unmarshal([]byte(MultiConfSettings), &multiconf.settings)

	// higher precedence
	os.Setenv("MULTICONF_COUNT", "15")
	os.Setenv("MULTICONF_ARGS", "arg-1,arg-2")

	// highest precedence
	args := []string{"--count", "12", "--args", "arg-3"}

	ops := []cmder.ExecuteOption{
		cmder.WithArgs(args),
		cmder.WithEnvironmentBinding(),
	}

	if err := cmder.Execute(context.Background(), multiconf, ops...); err != nil {
		fmt.Printf("unexpected error occurred: %v", err)
	}
	// Output:
	// format: pretty
	// count: 12
	// args: [arg-0 arg-1 arg-2 arg-3]
}

const MultiConfSettings = `{
	"format": "pretty",
	"count": 10,
	"args": ["arg-0"]
}
`

const MultiConfDesc = `
'multi-conf' desmonstrates how you can setup configuration from a configuration file (json), environment variables, and
command-line flags. In this example, configuration is evaluated in this order, from lowest to highest precedence:

  1. Configuration File  (/etc/multi.conf)
  2. Environment Variables (MULTICONF_*)
  3. Command Flags
`

const MultiConfExamples = `
$ MULTICONF_COUNT=15 MULTICONF_ARGS=arg-1,arg-2 multi-conf --count 12 --args arg-3
format: pretty
count: 12
args: [arg-0 arg-1 arg-2 arg-3]
`

type MultiConfConfig struct {
	Format string   `json:"format"`
	Count  int      `json:"count"`
	Args   []string `json:"args"`
}

type MultiConf struct {
	cmder.CommandDocumentation

	settings MultiConfConfig
}

var (
	multiconf = &MultiConf{
		CommandDocumentation: cmder.CommandDocumentation{
			Usage:     "multi-conf [flags]",
			ShortHelp: "Simple demonstration of application configuration.",
			Help:      MultiConfDesc,
			Examples:  MultiConfExamples,
		},
	}
)

func (m *MultiConf) Name() string {
	return "multi-conf"
}

func (m *MultiConf) InitializeFlags(fs *flag.FlagSet) {
	fs.StringVar(&m.settings.Format, "format", m.settings.Format, "specify a format")
	fs.IntVar(&m.settings.Count, "count", m.settings.Count, "specify a count")
	fs.Var((*getopt.StringsVar)(&m.settings.Args), "args", "provide arguments")
}

func (m *MultiConf) Run(ctx context.Context, args []string) error {
	fmt.Printf("format: %s\n", m.settings.Format)
	fmt.Printf("count: %d\n", m.settings.Count)
	fmt.Printf("args: %v\n", m.settings.Args)
	return nil
}
