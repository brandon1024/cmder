package getopt_test

import (
	"flag"
	"fmt"

	"github.com/brandon1024/cmder/getopt"
)

// This example demonstrates usage of [getopt.StringsVar] for string slice flags. You'll often find string sliceee flags
// on commands that accept IP addresses, for example.
func ExampleStringsVar() {
	hosts := getopt.StringsVar{}

	fs := flag.NewFlagSet("stringsvar", flag.ContinueOnError)
	fs.Var(&hosts, "broker", "connect to a broker")
	fs.Var(&hosts, "b", "connect to a broker")

	args := []string{
		"--broker", "tcp://127.0.0.1",
		"-b", "tls://broker-1.domain.example.com,tls://broker-2.domain.example.com",
	}

	if err := fs.Parse(args); err != nil {
		panic(err)
	}

	for _, host := range hosts {
		fmt.Printf("'%s'\n", host)
	}
	// Output:
	// 'tcp://127.0.0.1'
	// 'tls://broker-1.domain.example.com'
	// 'tls://broker-2.domain.example.com'
}
