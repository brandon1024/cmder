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

	// option 1: use MapVar directly
	variables := getopt.MapVar{}
	fs.Var(&variables, "variable", "specify runtime variables")

	// option 2: wrap an existing map with Map
	arg := map[string]string{}
	fs.Var(getopt.Map(arg), "arg", "specify runtime args")

	args := []string{
		"--variable", "key1=value1",
		"--variable", "key2=value2,key3=value3",
		`--arg="hello= HI, WORLD "`,
	}

	if err := fs.Parse(args); err != nil {
		panic(err)
	}

	for _, k := range slices.Sorted(maps.Keys(variables)) {
		fmt.Printf("%s: '%s'\n", k, variables[k])
	}
	for _, k := range slices.Sorted(maps.Keys(arg)) {
		fmt.Printf("%s: '%s'\n", k, arg[k])
	}
	// Output:
	// key1: 'value1'
	// key2: 'value2'
	// key3: 'value3'
	// hello: ' HI, WORLD '
}
