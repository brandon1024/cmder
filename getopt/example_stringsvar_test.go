package getopt_test

import (
	"flag"
	"fmt"

	"github.com/brandon1024/cmder/getopt"
)

// This example demonstrates usage of [getopt.StringsVar] for string slice flags. You'll often find string slice flags
// on commands that accept IP addresses, for example.
func ExampleStringsVar() {
	fs := flag.NewFlagSet("stringsvar", flag.ContinueOnError)

	// option 1: use StringsVar directly
	var hosts getopt.StringsVar
	fs.Var(&hosts, "broker", "connect to a broker")

	// option 2: wrap an existing slice
	var args []string
	fs.Var((*getopt.StringsVar)(&args), "a", "provide args")

	fs.Parse([]string{
		"--broker", "tls://broker-1.domain.example.com,tls://broker-2.domain.example.com",
		"-a", "CLIENT_USER",
		"-a", "CLIENT_PASS",
	})

	for _, host := range hosts {
		fmt.Printf("broker: '%s'\n", host)
	}
	for _, arg := range args {
		fmt.Printf("arg: '%s'\n", arg)
	}
	// Output:
	// broker: 'tls://broker-1.domain.example.com'
	// broker: 'tls://broker-2.domain.example.com'
	// arg: 'CLIENT_USER'
	// arg: 'CLIENT_PASS'
}
