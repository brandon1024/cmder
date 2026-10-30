package cmder

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/brandon1024/cmder/getopt"
)

// FlagInitializer is an interface implemented by a [Command] that need to register flags.
//
// InitializeFlags will be invoked during [Execute], prior to Initialize()/Run()/Destroy() routines. You can use this to
// register flags for your command.
//
// If the command does not define help flags '-h' or '--help', they will be registered automatically and will instruct
// [Execute] to render command usage.
type FlagInitializer interface {
	InitializeFlags(*flag.FlagSet)
}

// flagParser is an interface implemented by types that parse args.
type flagParser interface {
	Parse([]string) error
	Args() []string
}

// parseArgs processes args for the given command, returning the unparsed (remaining) arguments.
func (c command) parseArgs(args []string, ops *ExecuteOptions) ([]string, error) {
	var fp flagParser = &getopt.PosixFlagSet{
		FlagSet:        c.fs,
		RelaxedParsing: ops.relaxedFlags,
		Usage:          func() {},
	}

	if ops.nativeFlags {
		fp = c.fs
	}

	// interspersed args only possible for leaf commands
	interspersed := ops.interspersed
	if len(collectSubcommands(c.Command)) > 0 {
		interspersed = false
	}

	var processed []string

	for len(args) > 0 {
		if err := fp.Parse(args); err != nil {
			return nil, err
		}

		args = fp.Args()

		if !interspersed {
			return args, nil
		}

		if len(args) > 0 {
			processed = append(processed, args[0])
			args = args[1:]
		}
	}

	return processed, nil
}

// bindEnvironmentFlags sets flag values from matching environment variables.
func (c command) bindEnvironmentFlags(ops *ExecuteOptions) error {
	var components []string

	for _, parent := range c.parents {
		components = append(components, parent.Name())
	}

	components = append(components, c.Name())

	var flags []*flag.Flag
	c.fs.VisitAll(func(f *flag.Flag) {
		flags = append(flags, f)
	})

	for _, flag := range flags {
		variable := ops.bindEnvPrefix + formatEnvvar(append(components, flag.Name))

		if value, ok := os.LookupEnv(variable); ok {
			if err := flag.Value.Set(value); err != nil {
				return errors.Join(
					ErrEnvironmentBindFailure,
					fmt.Errorf("cmder: failed to set flag %s from variable %s", flag.Name, variable),
					err,
				)
			}
		}
	}

	return nil
}

// formatEnvvar generates an environment variable name which maps to the given flag path.
func formatEnvvar(flagPath []string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9]+")

	for i, v := range flagPath {
		flagPath[i] = strings.ToUpper(reg.ReplaceAllString(v, ""))
	}

	return strings.Join(flagPath, "_")
}
