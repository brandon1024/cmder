package getopt_test

import (
	"flag"
	"fmt"
	"time"

	"github.com/brandon1024/cmder/getopt"
)

// This example demonstrates the usage of [getopt.TimeVar].
func ExampleTimeVar() {
	fs := flag.NewFlagSet("custom", flag.ContinueOnError)

	var since, until time.Time

	fs.Var(getopt.Time(&since), "since", "show items since")
	fs.Var(getopt.Time(&until), "until", "show items until")

	args := []string{
		"-since", "2025-01-01T00:00:00Z",
		"-until", "2026-01-01T00:00:00Z",
	}

	if err := fs.Parse(args); err != nil {
		panic(err)
	}

	fmt.Printf("since: %s\n", since.String())
	fmt.Printf("until: %s\n", until.String())
	// Output:
	// since: 2025-01-01 00:00:00 +0000 UTC
	// until: 2026-01-01 00:00:00 +0000 UTC
}
