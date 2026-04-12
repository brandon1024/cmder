package getopt

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

// PosixFlagSet a wrapper over the standard [flag.FlagSet] that parses arguments with getopt-style (GNU/POSIX) semantics
// with short and long options.
//
// # Usage
//
// Start by initializing a new [PosixFlagSet].
//
//	fs := getopt.NewPosixFlagSet("hello", flag.ContinueOnError)
//
// You may also wrap an existing standard [flag.FlagSet] if you prefer.
//
//	fs := flag.NewFlagSet("hello", flag.ContinueOnError)
//	gfs := &getopt.PosixFlagSet{FlagSet: fs}
//
// Add flags to your PosixFlagSet with [PosixFlagSet.StringVar], [PosixFlagSet.BoolVar], [PosixFlagSet.IntVar], etc.
//
//	var (
//		output string
//		all    bool
//		count  int
//	)
//
//	fs.StringVar(&output, "output", "-", "output file location")
//	fs.BoolVar(&all, "a", false, "show all")
//	fs.IntVar(&count, "c", 0, "limit results to count")
//	fs.IntVar(&count, "count", 0, "limit results to count")
//
// The example above declares a long string flag '--output', a short '-a' bool flag, and aliased flags '-c' / '--count'.
//
// After all flags are defined, call [PosixFlagSet.Parse] to parse the flags.
//
//	err := fs.Parse(os.Args[1:])
//
// One parsed, any remaining (unparsed) arguments can be accessed with [PosixFlagSet.Arg] or [PosixFlagSet.Args].
//
// # Syntax
//
// [PosixFlagSet] distinguishes between long and short flags. A long flag is any flag whose name contains more than a
// single character, while a short flag has a name with a single character.
//
//	-a          // short boolean flag
//	--all       // long boolean flag
//	--all=false // disabled long boolean flag
//	-c 12       // short integer flag
//	-c12        // short integer flag with immediate value
//	--count 12  // long integer value
//	--count=12  // long integer value with immediate value
//
// Short boolean flags may be combined into a single argument, and short flags accepting arguments may be "stuck" to the
// value:
//
//	-ac12       // equivalent to '-a -c 12'
//
// Flag parsing stops just before the first non-flag argument ("-" is a non-flag argument) or after the terminator "--".
//
// Flags which accept a number ([PosixFlagSet.Int], [PosixFlagSet.Uint], [PosixFlagSet.Float64], etc) will parse their arguments with
// [strconv]. For integers, binary/octal/decimal/hexadecimal numbers are accepted (see [strconv.ParseInt] and
// [strconv.ParseUint]). For floats, anything parseable by [strconv.ParseFloat] is accepted.
//
//	--count 12
//	--count 0xC
//	--count 0o14
//	--count 0b1100
//	--count 1.2E1
//
// Boolean flags with an immediate value may be anything parseable by [strconv.ParseBool].
//
//	--all=false
//	--all=FALSE
//	--all=f
//	--all=0
//
// Duration flags accept any input valid for [time.ParseDuration].
//
//	--since=3m2s
type PosixFlagSet struct {
	*flag.FlagSet

	// Similar to [flag.FlagSet.Usage], Usage is invoked when parsing fails. By default, uses
	// [PosixFlagSet.PrintDefaults] which renders flag usage with posix semantics.
	Usage func()

	// If true, relaxes flag parsing allowing Parse to accept partial flag matches (e.g. '--auto' for '--auto-gc'). An
	// error will still be emitted if the input is ambiguous (e.g. '--auto' for '--auto-gc' or '--auto-maintenance').
	RelaxedParsing bool

	parsed bool
	args   []string
}

// NewPosixFlagSet builds a new [flag.FlagSet] and wraps it with a [PosixFlagSet].
func NewPosixFlagSet(name string, e flag.ErrorHandling) *PosixFlagSet {
	return &PosixFlagSet{
		FlagSet: flag.NewFlagSet(name, e),
	}
}

// PrintDefaults prints usage information and default values for all flags of this flag set to the output location
// configured with [flag.FlagSet.Init] or [flag.FlagSet.SetOutput].
func (f *PosixFlagSet) PrintDefaults() {
	f.VisitAll(func(flg *flag.Flag) {
		var err error

		if isHiddenFlag(flg) {
			return
		}

		name, usage := flag.UnquoteUsage(flg)

		if len(flg.Name) == 1 {
			_, err = fmt.Fprintf(f.Output(), "   -%s", flg.Name)
		} else {
			_, err = fmt.Fprintf(f.Output(), "  --%s", flg.Name)
		}

		if err != nil {
			panic(err)
		}

		if len(name) > 0 && !isBoolFlag(flg) {
			_, err = fmt.Fprintf(f.Output(), " <%s>", name)
		}

		if err != nil {
			panic(err)
		}

		if len(flg.DefValue) > 0 {
			_, err = fmt.Fprintf(f.Output(), " (default %s)", flg.DefValue)
		}

		if err != nil {
			panic(err)
		}

		_, err = fmt.Fprintf(f.Output(), "\n        %s\n", usage)

		if err != nil {
			panic(err)
		}
	})
}

// Arg returns the i'th remaining argument after calling [PosixFlagSet.Parse]. Returns an empty string if the argument does
// not exist, or [PosixFlagSet.Parse] was not called.
func (f *PosixFlagSet) Arg(i int) string {
	if i < 0 || i >= len(f.args) {
		return ""
	}

	return f.args[i]
}

// NArg returns the number of non-flag arguments remaining after calling [PosixFlagSet.Parse].
func (f *PosixFlagSet) NArg() int {
	return len(f.args)
}

// Args returns a slice of non-flag arguments remaining after calling [PosixFlagSet.Parse].
func (f *PosixFlagSet) Args() []string {
	return f.args
}

// Parsed returns whether or not [PosixFlagSet.Parse] has been invoked on this flag set.
func (f *PosixFlagSet) Parsed() bool {
	return f.parsed
}

// Parse processes the given arguments and updates the flags of this flag set. The arguments given should not include
// the command name. Parse should only be called after all flags have been registered and before flags are accessed by
// the application.
//
// The return value will be [flag.ErrHelp] if -help or -h were set but not defined.
func (f *PosixFlagSet) Parse(arguments []string) error {
	usage := f.Usage
	if usage == nil {
		usage = f.PrintDefaults
	}

	err := f.parse(arguments)
	if err == nil {
		return nil
	}

	if f.ErrorHandling() == flag.ContinueOnError {
		usage()
		return err
	}

	if f.ErrorHandling() == flag.PanicOnError {
		usage()
		panic(err)
	}

	if errors.Is(err, flag.ErrHelp) && f.ErrorHandling() == flag.ExitOnError {
		usage()
		os.Exit(0)
		return nil
	}

	usage()
	os.Exit(2)
	return nil
}

func (f *PosixFlagSet) parse(arguments []string) error {
	var err error

	f.parsed = true

	for len(arguments) > 0 {
		arg := arguments[0]

		// a single hyphen is not a flag -- update arguments and return
		if arg == "-" {
			f.args = arguments[0:]
			return nil
		}

		// double hyphens is sentinel and denotes end of arguments -- remove from arguments and return
		if arg == "--" {
			f.args = arguments[1:]
			return nil
		}

		// parse long option
		long, ok := strings.CutPrefix(arg, "--")
		if ok {
			arguments, err = f.parseLong(long, arguments[1:])
			if err != nil {
				return err
			}

			continue
		}

		// parse short option
		short, ok := strings.CutPrefix(arg, "-")
		if ok {
			arguments, err = f.parseShort(short, arguments[1:])
			if err != nil {
				return err
			}

			continue
		}

		f.args = arguments
		return nil
	}

	f.args = arguments
	return nil
}

func (f *PosixFlagSet) parseLong(arg string, arguments []string) ([]string, error) {
	arg, value, inlineVal := strings.Cut(arg, "=")

	flg := f.lookupLong(arg, f.RelaxedParsing)

	// similar to the stdlib, if we encounter a '--help' flag but none defined, return ErrHelp
	if flg == nil && arg == "help" {
		return nil, flag.ErrHelp
	}

	if flg == nil {
		return nil, fmt.Errorf("flag '--%s' does not exist", arg)
	}

	if isBoolFlag(flg) {
		if !inlineVal {
			value = "true"
		}
	} else {
		// if the value was not provided inline '--arg=value', grab the next argument
		if !inlineVal {
			if len(arguments) == 0 {
				return nil, fmt.Errorf("missing argument to flag '--%s'", arg)
			}

			arguments, value = arguments[1:], arguments[0]
		}
	}

	if err := f.Set(flg.Name, value); err != nil {
		return nil, err
	}

	return arguments, nil
}

func (f *PosixFlagSet) parseShort(short string, arguments []string) ([]string, error) {
	for len(short) > 0 {
		args := strings.SplitN(short, "", 2)

		if len(args) == 1 {
			short = ""
		}
		if len(args) == 2 {
			short = args[1]
		}

		flg := f.Lookup(args[0])
		if flg == nil && args[0] == "h" {
			return nil, flag.ErrHelp
		}
		if flg == nil {
			return nil, fmt.Errorf("flag '-%s' does not exist", args[0])
		}

		if isBoolFlag(flg) {
			if err := f.Set(args[0], "true"); err != nil {
				return nil, err
			}
		} else {
			if short != "" {
				// rest is arg
				if err := f.Set(args[0], short); err != nil {
					return nil, err
				}
			} else {
				// take next arg
				if len(arguments) == 0 {
					return nil, fmt.Errorf("missing argument to flag '-%s'", args[0])
				}

				if err := f.Set(args[0], arguments[0]); err != nil {
					return nil, err
				}

				arguments = arguments[1:]
			}

			return arguments, nil
		}
	}

	return arguments, nil
}

// lookupLong looks for a (long) flag with the given name in f. Returns nil if no flag found.
//
// When relaxed is true, partial flag name matches are permitted. If more than one flag name has the prefix name,
// returns nil.
func (f *PosixFlagSet) lookupLong(name string, relaxed bool) *flag.Flag {
	// never match a short name, since the user is expected to use short-style flags instead ('-a' and not '--a')
	if len(name) <= 1 {
		return nil
	}

	var flags []*flag.Flag

	f.VisitAll(func(flg *flag.Flag) {
		// don't match short flags
		if len(flg.Name) <= 1 {
			return
		}

		if !relaxed && flg.Name == name {
			flags = append(flags, flg)
		}
		if relaxed && strings.HasPrefix(flg.Name, name) {
			flags = append(flags, flg)
		}
	})

	if len(flags) != 1 {
		return nil
	}

	return flags[0]
}
