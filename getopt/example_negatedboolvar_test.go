package getopt_test

import (
	"flag"
	"fmt"
	"os"

	"github.com/brandon1024/cmder/getopt"
)

// This example demonstrates the usage of [getopt.NegatedBoolVar].
func ExampleNegatedBoolVar() {
	var (
		sign   bool
		verify bool
	)

	fs := getopt.NewPosixFlagSet("custom", flag.ContinueOnError)

	// option 1: using NegatedBoolVar directly
	fs.BoolVar(&sign, "gpg-sign", false, "gpg sign the input")
	fs.Var((*getopt.NegatedBoolVar)(&sign), "no-gpg-sign", "skip gpg signing")

	// option 2: with NegatedBool
	fs.BoolVar(&verify, "verify", false, "verify the result")
	fs.Var(getopt.NegatedBool(&verify), "no-verify", "skip result verification")

	fs.SetOutput(os.Stdout)
	fs.PrintDefaults()

	args := []string{
		"--gpg-sign", "--no-gpg-sign",
		"--no-verify=false",
	}

	if err := fs.Parse(args); err != nil {
		panic(err)
	}

	fmt.Printf("sign: %v\n", sign)
	fmt.Printf("verify: %v\n", verify)
	// Output:
	//   --gpg-sign, --no-gpg-sign
	//       gpg sign the input
	//
	//   --verify, --no-verify
	//       verify the result
	// sign: false
	// verify: true
}
