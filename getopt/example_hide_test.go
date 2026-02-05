package getopt_test

import (
	"flag"
	"fmt"
	"os"

	"github.com/brandon1024/cmder/getopt"
)

// This example demonstrates the usage of [getopt.Hidden] and [getopt.Hide].
func ExampleHide() {
	var (
		count  uint
		output string
	)

	fs := getopt.NewPosixFlagSet("hidden", flag.ContinueOnError)

	fs.UintVar(&count, "count", 12, "`number` of results")
	fs.UintVar(&count, "c", 12, "`number` of results (hidden flag)")
	fs.StringVar(&output, "output", "-", "output `file`")
	fs.StringVar(&output, "o", "-", "output `file`")

	// option 1: wrap the flag value
	fs.Lookup("c").Value = &getopt.Hidden{fs.Lookup("c").Value}

	// option 2: use getopt.Hide
	getopt.Hide(fs.Lookup("o"))

	// option 3: using FlagSet.Var
	var since getopt.TimeVar
	fs.Var(&getopt.Hidden{&since}, "since", "show items since")

	fs.SetOutput(os.Stdout)
	fs.PrintDefaults()

	args := []string{"-c", "2025", "-o", "output.txt", "--since", "2025-01-01T00:00:00Z"}

	if err := fs.Parse(args); err != nil {
		panic(err)
	}

	fmt.Printf("values: %d %s %s\n", count, output, since.String())

	// Output:
	//   --count <number> (default 12)
	//         number of results
	//   --output <file> (default -)
	//         output file
	// values: 2025 output.txt 2025-01-01T00:00:00Z
}
