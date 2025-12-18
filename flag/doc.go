/*
Package flag is an alternative to the standard library package of the same name. This implementation adheres to
POSIX/GNU standards for option and argument parsing. This package aims to be a simple drop-in replacement for the
standard implementation and will continue to track the standard library flag package. No new features deviating from the
standard implementation will be accepted.

Unless otherwise mentioned, the documentation for the standard library implementation applies here too. Refer to the
documentation there for more information:

	https://pkg.go.dev/flag

# Usage

Start by initializing a new [FlagSet].

	fs := flag.NewFlagSet("hello", ContinueOnError)

Add flags to your FlagSet with [FlagSet.StringVar], [FlagSet.BoolVar], [FlagSet.IntVar], etc.

	var (
		output string
		all    bool
		count  int
	)

	fs.StringVar(&output, "output", "-", "output file location")
	fs.BoolVar(&all, "a", false, "show all")
	fs.IntVar(&count, "c", 0, "limit results to count")
	fs.IntVar(&count, "count", 0, "limit results to count")

The example above declares a long string flag '--output', a short '-a' bool flag, and aliased flags '-c' / '--count'.

After all flags are defined, call [FlagSet.Parse] to parse the flags.

	err := fs.Parse(os.Args[1:])

One parsed, any remaining (unparsed) arguments can be accessed with [FlagSet.Arg] or [FlagSet.Args].

# Syntax

[FlagSet] distinguishes between long and short flags. A long flag is any flag whose name contains more than a single
character, while a short flag has a name with a single character.

	-a          // short boolean flag
	--all       // long boolean flag
	--all=false // disabled long boolean flag
	-c 12       // short integer flag
	-c12        // short integer flag with immediate value
	--count 12  // long integer value
	--count=12  // long integer value with immediate value

Multiple short flags may be "stuck" together.

	-ac12       // equivalent to '-a -c 12'

Flag parsing stops just before the first non-flag argument ("-" is a non-flag argument) or after the terminator "--".

Flags which accept a number ([FlagSet.Int], [FlagSet.Uint], [FlagSet.Float64], etc) will parse their arguments with
[strconv]. For integers, binary/octal/decimal/hexadecimal numbers are accepted (see [strconv.ParseInt] and
[strconv.ParseUint]). For floats, anything parseable by [strconv.ParseFloat] is accepted.

	--count 12
	--count 0xC
	--count 0o14
	--count 0b1100
	--count 1.2E1

Boolean flags with an immediate value may be anything parseable by [strconv.ParseBool].

	--all=false
	--all=FALSE
	--all=f
	--all=0

Duration flags accept any input valid for [time.ParseDuration].

	--since=3m2s
*/
package flag
