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

	var hosts, args, patterns []string

	fs.Var(getopt.Strings(&hosts), "broker", "connect to a broker")
	fs.Var(getopt.Strings(&args), "a", "provide args")
	fs.Var(getopt.Strings(&patterns), "p", "provide patterns")

	fs.Parse([]string{
		"--broker", "tls://broker-1.domain.example.com,tls://broker-2.domain.example.com",
		"-a", "CLIENT_USER",
		"-a", "CLIENT_PASS",
		"-p", "**/*.go,*.mod,*.sum",
	})

	for _, host := range hosts {
		fmt.Printf("broker: '%s'\n", host)
	}
	for _, arg := range args {
		fmt.Printf("arg: '%s'\n", arg)
	}
	for _, pattern := range patterns {
		fmt.Printf("patterns: '%s'\n", pattern)
	}
	// Output:
	// broker: 'tls://broker-1.domain.example.com'
	// broker: 'tls://broker-2.domain.example.com'
	// arg: 'CLIENT_USER'
	// arg: 'CLIENT_PASS'
	// patterns: '**/*.go'
	// patterns: '*.mod'
	// patterns: '*.sum'
}
