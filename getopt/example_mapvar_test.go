package getopt_test

import (
	"flag"
	"fmt"
	"maps"
	"slices"

	"github.com/brandon1024/cmder/getopt"
)

// This example demonstrates usage of [getopt.MapVar] for string maps. You'll often find map flags on commands that
// perform templating of text files, for example.
func ExampleMapVar() {
	fs := flag.NewFlagSet("map", flag.ContinueOnError)

	arg := map[string]string{}
	fs.Var(getopt.Map(arg), "arg", "specify runtime args")

	args := []string{
		"--arg", "key1=value1",
		"--arg", "key2=value2,key3=value3",
		`--arg="hello= HI, WORLD "`,
	}

	if err := fs.Parse(args); err != nil {
		panic(err)
	}

	for _, k := range slices.Sorted(maps.Keys(arg)) {
		fmt.Printf("%s: '%s'\n", k, arg[k])
	}
	// Output:
	// hello: ' HI, WORLD '
	// key1: 'value1'
	// key2: 'value2'
	// key3: 'value3'
}
