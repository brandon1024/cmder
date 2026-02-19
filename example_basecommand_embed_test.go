package cmder_test

import (
	"context"
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/brandon1024/cmder"
	"github.com/brandon1024/cmder/getopt"
)

// This example demonstrates an alternative usage of [cmder.BaseCommand]. By embedding BaseCommand into your own types,
// your type implements all required interfaces needed to fulfill [cmder.Command]. You can override the standard methods
// with your own implementation.
func ExampleBaseCommand_embedding() {
	args := []string{"-m", "1.6.17-beta.0+20130313144700"}

	if err := cmder.Execute(context.Background(), semver, cmder.WithArgs(args)); err != nil {
		fmt.Printf("unexpected error occurred: %v", err)
	}
	// Output:
	// 1.7.0-beta.0
}

const SemverDesc = `
'semver' demonstrates an alternative usage of [cmder.BaseCommand]. By embedding BaseCommand into your own types, your
type implements all required interfaces needed to fulfill [cmder.Command]. You can override the standard methods with
your own implementation.

Example 'semver' manipulates a version string, bumping it up a major/minor/patch version or attaches build/pre-release
information. By default, input versoins are bumped up a patch version.
`

const SemverExamples = `
semver 0.5.3
semver --minor 1.6.17-alpha
semver --pre-release alpha --build ${BUILD_ID} 1.3.4
`

var (
	semver = &Semver{
		BaseCommand: cmder.BaseCommand{
			CommandName: "semver",
			CommandDocumentation: cmder.CommandDocumentation{
				Usage:     "semver [--major | --minor | --patch] [--pre-release <pre>] [--build <build>] <version>",
				ShortHelp: "A simple demonstration of embedding BaseCommand in custom types.",
				Help:      SemverDesc,
				Examples:  SemverExamples,
			},
		},
	}
)

type Semver struct {
	cmder.BaseCommand

	major, minor, patch bool
	pre, build          string
}

// Configure command flags. Register short aliases for all long flag options.
//
// InitializeFlags overrides the default [cmder.BaseCommand.InitializeFlags] implementation.
func (s *Semver) InitializeFlags(fs *flag.FlagSet) {
	fs.BoolVar(&s.major, "major", s.major, "bump version to the next major version")
	fs.BoolVar(&s.minor, "minor", s.minor, "bump version to the next minor version")
	fs.BoolVar(&s.patch, "patch", s.patch, "bump version to the next patch version")
	fs.StringVar(&s.pre, "pre-release", s.pre, "include pre-release information in output (e.g. alpha, x.7.z.92)")
	fs.StringVar(&s.build, "build", s.build, "include build information in output (e.g. 20130313144700, exp.sha.5114f85)")

	getopt.Alias(fs, "major", "M")
	getopt.Alias(fs, "minor", "m")
	getopt.Alias(fs, "patch", "p")
	getopt.Alias(fs, "pre-release", "x")
	getopt.Alias(fs, "build", "b")
}

// The command's run function.
//
// Run overrides the default [cmder.BaseCommand.Run] implementation.
func (s *Semver) Run(ctx context.Context, args []string) error {
	if len(args) != 1 {
		return cmder.ErrShowUsage
	}

	version, _, _ := strings.Cut(args[0], "+")
	version, pre, _ := strings.Cut(version, "-")

	var (
		parts []string
		level = 2
	)

	if s.minor {
		level = 1
	}
	if s.major {
		level = 0
	}

	for i, v := range strings.Split(version, ".") {
		num, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("semver: invalid input: %s", args[0])
		}
		if i == level {
			num++
		}
		if i > level {
			num = 0
		}

		parts = append(parts, strconv.Itoa(num))
	}

	version = strings.Join(parts, ".")

	if s.pre != "" {
		version += "-" + s.pre
	} else if pre != "" {
		version += "-" + pre
	}

	if s.build != "" {
		version += "+" + s.build
	}

	fmt.Println(version)

	return nil
}
