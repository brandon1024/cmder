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
	variables := getopt.MapVar{}

	fs := flag.NewFlagSet("map", flag.ContinueOnError)
	fs.Var(&variables, "variable", "specify runtime variables")
	fs.Var(&variables, "v", "specify runtime variables")

	args := []string{
		"--variable", "key1=value1",
		"-v", "key2=value2,key3=value3",
		`--variable="hello= HI, WORLD "`,
	}

	if err := fs.Parse(args); err != nil {
		panic(err)
	}

	for _, k := range slices.Sorted(maps.Keys(variables)) {
		fmt.Printf("%s: '%s'\n", k, variables[k])
	}
	// Output:
	// hello: ' HI, WORLD '
	// key1: 'value1'
	// key2: 'value2'
	// key3: 'value3'
}
