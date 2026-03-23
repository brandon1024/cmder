package getopt_test

import (
	"flag"
	"fmt"

	"github.com/brandon1024/cmder/getopt"
)

// This example demonstrates the usage of [getopt.Counter].
func ExampleCounter() {
	var verbose uint

	fs := getopt.NewPosixFlagSet("counter", flag.ContinueOnError)
	fs.Var(getopt.Counter(&verbose), "v", "increase verbosity")

	if err := fs.Parse([]string{"-vvv"}); err != nil {
		panic(err)
	}

	fmt.Printf("verbosity: %d\n", verbose)
	// Output:
	// verbosity: 3
}
