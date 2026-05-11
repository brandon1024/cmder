package getopt

import (
	"cmp"
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"slices"
	"strings"
	"text/template"
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

// PrintDefaults writes usage information and default values for all flags in the flag set to the output configured by
// [flag.FlagSet.Init] or [flag.FlagSet.SetOutput].
//
// Unlike [flag.FlagSet.PrintDefaults] in the standard library, this method formats flags using GNU/POSIX conventions,
// such as '-a' for short options and '--all' for long options.
//
// Flags are grouped by [flag.Value] equivalence so that aliases appear together in the usage text. For example, a short
// flag may be an alias of a long flag (see [Alias]):
//
//	-a <string>, --addr=<string>
//	-s <string>, --serial-number=<string>
//
// Hidden flags, created with [Hide], are omitted from the output.
func (f *PosixFlagSet) PrintDefaults() {
	format := `
		{{- $print_started := false -}}

		{{- range . -}}
			{{- if $print_started -}}
				{{- println -}}
			{{- end -}}
			{{- $print_started = true -}}

			{{- printf "  " -}}

			{{- range $index, $flg := . -}}
				{{- if (ne $index 0) -}}
					{{- printf ", " -}}
				{{- end -}}

				{{- if (eq (len $flg.Name) 1) -}}
					{{- printf "-%s" .Name -}}
				{{- else -}}
					{{- printf "--%s" .Name -}}
				{{- end -}}

				{{- $name := (index (unquote $flg) 0) -}}

				{{- if (bool $flg) -}}
				{{- else if (and $name (eq (len $flg.Name) 1)) -}}
					{{- printf " <%s>" $name -}}
				{{- else if $name -}}
					{{- printf "=<%s>" $name -}}
				{{- end -}}
			{{- end -}}

			{{ if (not (zero (index . 0))) }}
				{{- printf " (default %s)" (index . 0).DefValue -}}
			{{- end -}}

			{{- println -}}

			{{- printf "      %s\n" (index (unquote (index . 0)) 1) -}}
		{{- end -}}`

	tmpl, err := template.New("usage").Funcs(template.FuncMap{
		"unquote": unquote,
		"zero":    zero,
		"bool":    isBoolFlag,
	}).Parse(format)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(f.Output(), f.group())
	if err != nil {
		panic(err)
	}
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
		usage = f.defaultUsage
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

// group organizes the flags and returns them.
//
// The flags are grouped by [flag.Value] equivalence. This allows flags to be grouped together in the rendered
// usage text when two flags are aliases of each other. This is often the case for short flags which are aliases of
// longer flags (e.g. '-a' is an alias of '--all').
//
//	-a <string>, --addr=<string>
//	-s <string>, --serial-number=<string>
//
// The resulting map entries are keyed by the flag group name, which is the longest flag name in the group. The map
// values are slices of (one or more) flags in the flag group, sorted by flag name length ('-a' before '--all').
//
// Hidden flags are excluded from the resulting map.
func (f *PosixFlagSet) group() map[string][]*flag.Flag {
	var collected []*flag.Flag

	f.VisitAll(func(f *flag.Flag) {
		if !isHiddenFlag(f) {
			collected = append(collected, f)
		}
	})

	// sort flags by name length in descending order to ensure that keys in resulting map will use long names first
	slices.SortFunc(collected, func(a, b *flag.Flag) int {
		return cmp.Compare(len(b.Name), len(a.Name))
	})

	groups := map[string][]*flag.Flag{}

	for len(collected) > 0 {
		var flg *flag.Flag

		// pop the head of the slice
		flg, collected = collected[0], collected[1:]

		// update groups
		groups[flg.Name] = []*flag.Flag{flg}

		// traverse the flags again and find (and remove) any which match flg
		for i := len(collected) - 1; i >= 0; i-- {
			other := collected[i]

			if areSame(flg.Value, other.Value) {
				groups[flg.Name] = append(groups[flg.Name], other)
				collected = append(collected[:i], collected[i+1:]...)
			}
		}

		// sort by length (then lexical order), this time ascending (-a before --all)
		slices.SortFunc(groups[flg.Name], func(a, b *flag.Flag) int {
			if c := cmp.Compare(len(a.Name), len(b.Name)); c != 0 {
				return c
			}

			return cmp.Compare(a.Name, b.Name)
		})
	}

	return groups
}

// defaultUsage is the default usage renderer invoked when parsing fails, invoked when [PosixFlagSet.Usage] is nil.
func (f *PosixFlagSet) defaultUsage() {
	name := f.Name()

	if name == "" {
		_, _ = fmt.Fprintf(f.Output(), "Usage:\n")
	} else {
		_, _ = fmt.Fprintf(f.Output(), "Usage of %s:\n", name)
	}

	f.PrintDefaults()
}

// unquote is a wrapper over the standard [flag.UnquoteUsage] which returns a slice, allowing it to be used as a
// template func.
func unquote(flg *flag.Flag) []string {
	name, usage := flag.UnquoteUsage(flg)
	return []string{name, usage}
}

// zero checks if the default value of flg is the zero value for its type. This is used when rendering usage text
// to render default flag values only when the default value is interesting.
//
// This function expects that flg adhere's to the same requirements of the stdlib [flag] package, notably:
//
//	The flag package may call the String method with a zero-valued receiver, such as a nil pointer.
//
// Flags that don't respect this requirement will result in an error.
func zero(flg *flag.Flag) (ok bool, err error) {
	var z reflect.Value

	if typ := reflect.TypeOf(flg.Value); typ.Kind() == reflect.Pointer {
		z = reflect.New(typ.Elem())
	} else {
		z = reflect.Zero(typ)
	}

	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("cmder: flag '%s' is backed by a type that does not accept calling String() on the zero value (bug): %v",
				flg.Name, e)
		}
	}()

	ok = flg.DefValue == z.Interface().(flag.Value).String()
	return
}

// areSame check if f1 and f2 have the same underlying [flag.Value].
func areSame(f1, f2 flag.Value) bool {
	var (
		ref1 = reflect.ValueOf(f1)
		ref2 = reflect.ValueOf(f2)
	)

	if ref1.Comparable() && ref2.Comparable() && f1 == f2 {
		return true
	}

	if ref1.Kind() != ref2.Kind() {
		return false
	}

	if !slices.Contains([]reflect.Kind{reflect.Map, reflect.Pointer, reflect.Func, reflect.Slice}, ref1.Kind()) {
		return false
	}

	return ref1.Pointer() == ref2.Pointer()
}
