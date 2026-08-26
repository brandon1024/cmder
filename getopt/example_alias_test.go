package getopt_test

import (
	"flag"
	"fmt"
	"os"

	"github.com/brandon1024/cmder/getopt"
)

// This example demonstrates the usage of [getopt.Alias].
func ExampleAlias() {
	var (
		count  uint
		output string
	)

	fs := getopt.NewPosixFlagSet("alias", flag.ContinueOnError)

	fs.UintVar(&count, "count", 12, "`number` of results")
	fs.StringVar(&output, "output", "-", "output `file`")

	getopt.Alias(fs.FlagSet, "count", "c")
	getopt.Alias(fs.FlagSet, "output", "o")

	fs.SetOutput(os.Stdout)
	fs.PrintDefaults()

	if err := fs.Parse([]string{"-c", "2025", "-o", "output.txt"}); err != nil {
		panic(err)
	}

	fmt.Printf("values: %d %s\n", count, output)

	// Output:
	//   -c <number>, --count=<number> (default 12)
	//       number of results
	//
	//   -o <file>, --output=<file> (default -)
	//       output file
	// values: 2025 output.txt
}
