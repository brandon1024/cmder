package getopt_test

import (
	"flag"
	"fmt"

	"github.com/brandon1024/cmder/getopt"
)

// This example demonstrates the usage of [getopt.TimeVar].
func ExampleTimeVar() {
	var since getopt.TimeVar

	fs := flag.NewFlagSet("custom", flag.ContinueOnError)
	fs.Var(&since, "since", "show items since")

	args := []string{
		"-since", "2025-01-01T00:00:00Z",
	}

	if err := fs.Parse(args); err != nil {
		panic(err)
	}

	fmt.Printf("since: %s\n", since.String())
	// Output:
	// since: 2025-01-01T00:00:00Z
}
