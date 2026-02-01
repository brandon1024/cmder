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

	fs.Lookup("c").Value = &getopt.Hidden{fs.Lookup("c").Value}
	getopt.Hide(fs.Lookup("o"))

	fs.SetOutput(os.Stdout)
	fs.PrintDefaults()

	args := []string{"-c", "2025", "-o", "output.txt"}

	if err := fs.Parse(args); err != nil {
		panic(err)
	}

	fmt.Printf("values: %d %s\n", count, output)

	// Output:
	//   --count <number> (default 12)
	//         number of results
	//   --output <file> (default -)
	//         output file
	// values: 2025 output.txt
}
